package linux_backend

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/vito/garden/backend"
	"github.com/vito/garden/backend/linux_backend/bandwidth_manager"
	"github.com/vito/garden/backend/linux_backend/cgroups_manager"
	"github.com/vito/garden/backend/linux_backend/job_tracker"
	"github.com/vito/garden/backend/linux_backend/network"
	"github.com/vito/garden/backend/linux_backend/quota_manager"
	"github.com/vito/garden/command_runner"
)

type LinuxContainer struct {
	id   string
	path string

	spec backend.ContainerSpec

	resources Resources

	portPool PortPool

	runner command_runner.CommandRunner

	cgroupsManager   cgroups_manager.CgroupsManager
	quotaManager     quota_manager.QuotaManager
	bandwidthManager bandwidth_manager.BandwidthManager

	jobTracker *job_tracker.JobTracker

	oomLock     sync.RWMutex
	oomNotifier *exec.Cmd
}

type Resources struct {
	UID     uint32
	Network network.Network
}

type PortPool interface {
	Acquire() (uint32, error)
	Release(uint32)
}

func NewLinuxContainer(
	id, path string,
	spec backend.ContainerSpec,
	resources Resources,
	portPool PortPool,
	runner command_runner.CommandRunner,
	cgroupsManager cgroups_manager.CgroupsManager,
	quotaManager quota_manager.QuotaManager,
	bandwidthManager bandwidth_manager.BandwidthManager,
) *LinuxContainer {
	return &LinuxContainer{
		id:   id,
		path: path,

		spec: spec,

		resources: resources,

		portPool: portPool,

		runner: runner,

		cgroupsManager:   cgroupsManager,
		quotaManager:     quotaManager,
		bandwidthManager: bandwidthManager,

		jobTracker: job_tracker.New(path, runner),
	}
}

func (c *LinuxContainer) ID() string {
	return c.id
}

func (c *LinuxContainer) Handle() string {
	if c.spec.Handle != "" {
		return c.spec.Handle
	}

	return c.ID()
}

func (c *LinuxContainer) Start() error {
	log.Println(c.id, "starting")

	start := exec.Command(path.Join(c.path, "start.sh"))

	start.Env = []string{
		"id=" + c.id,
		"container_iface_mtu=1500",
		"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
	}

	return c.runner.Run(start)
}

func (c *LinuxContainer) Stop(kill bool) error {
	log.Println(c.id, "stopping")

	stop := exec.Command(path.Join(c.path, "stop.sh"))

	if kill {
		stop.Args = append(stop.Args, "-w", "0")
	}

	err := c.runner.Run(stop)
	if err != nil {
		return err
	}

	c.stopOomNotifier()

	return nil
}

func (c *LinuxContainer) Info() (backend.ContainerInfo, error) {
	memoryStat, err := c.cgroupsManager.Get("memory", "memory.stat")
	if err != nil {
		return backend.ContainerInfo{}, err
	}

	cpuUsage, err := c.cgroupsManager.Get("cpuacct", "cpuacct.usage")
	if err != nil {
		return backend.ContainerInfo{}, err
	}

	cpuStat, err := c.cgroupsManager.Get("cpuacct", "cpuacct.stat")
	if err != nil {
		return backend.ContainerInfo{}, err
	}

	diskStat, err := c.quotaManager.GetUsage(c.resources.UID)
	if err != nil {
		return backend.ContainerInfo{}, err
	}

	bandwidthStat, err := c.bandwidthManager.GetLimits()
	if err != nil {
		return backend.ContainerInfo{}, err
	}

	return backend.ContainerInfo{
		State:         "active",   // TODO
		Events:        []string{}, // TODO
		HostIP:        c.resources.Network.HostIP().String(),
		ContainerIP:   c.resources.Network.ContainerIP().String(),
		ContainerPath: c.path,
		JobIDs:        c.jobTracker.ActiveJobs(),
		MemoryStat:    parseMemoryStat(memoryStat),
		CPUStat:       parseCPUStat(cpuUsage, cpuStat),
		DiskStat:      diskStat,
		BandwidthStat: bandwidthStat,
	}, nil
}

func (c *LinuxContainer) CopyIn(src, dst string) error {
	log.Println(c.id, "copying in from", src, "to", dst)
	return c.rsync(src, "vcap@container:"+dst)
}

func (c *LinuxContainer) CopyOut(src, dst, owner string) error {
	log.Println(c.id, "copying out from", src, "to", dst)

	err := c.rsync("vcap@container:"+src, dst)
	if err != nil {
		return err
	}

	if owner != "" {
		chown := exec.Command("chown", "-R", owner, dst)

		err := c.runner.Run(chown)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *LinuxContainer) LimitBandwidth(limits backend.BandwidthLimits) (backend.BandwidthLimits, error) {
	log.Println(
		c.id,
		"limiting bandwidth to",
		limits.RateInBytesPerSecond,
		"bytes per second; burst",
		limits.BurstRateInBytesPerSecond,
	)

	return limits, c.bandwidthManager.SetLimits(limits)
}

func (c *LinuxContainer) LimitDisk(limits backend.DiskLimits) (backend.DiskLimits, error) {
	err := c.quotaManager.SetLimits(c.resources.UID, limits)
	if err != nil {
		return backend.DiskLimits{}, err
	}

	return c.quotaManager.GetLimits(c.resources.UID)
}

func (c *LinuxContainer) LimitMemory(limits backend.MemoryLimits) (backend.MemoryLimits, error) {
	log.Println(c.id, "limiting memory to", limits.LimitInBytes, "bytes")

	err := c.startOomNotifier()
	if err != nil {
		return backend.MemoryLimits{}, err
	}

	limit := fmt.Sprintf("%d", limits.LimitInBytes)

	// memory.memsw.limit_in_bytes must be >= memory.limit_in_bytes
	//
	// however, it must be set after memory.limit_in_bytes, and if we're
	// increasing the limit, writing memory.limit_in_bytes first will fail.
	//
	// so, write memory.limit_in_bytes before and after
	c.cgroupsManager.Set("memory", "memory.limit_in_bytes", limit)

	err = c.cgroupsManager.Set("memory", "memory.memsw.limit_in_bytes", limit)
	if err != nil {
		return backend.MemoryLimits{}, err
	}

	err = c.cgroupsManager.Set("memory", "memory.limit_in_bytes", limit)
	if err != nil {
		return backend.MemoryLimits{}, err
	}

	actualLimitInBytes, err := c.cgroupsManager.Get("memory", "memory.limit_in_bytes")

	numericLimit, err := strconv.ParseUint(actualLimitInBytes, 10, 0)
	if err != nil {
		return backend.MemoryLimits{}, err
	}

	return backend.MemoryLimits{uint64(numericLimit)}, nil
}

func (c *LinuxContainer) Spawn(spec backend.JobSpec) (uint32, error) {
	log.Println(c.id, "spawning job:", spec.Script)

	wshPath := path.Join(c.path, "bin", "wsh")
	sockPath := path.Join(c.path, "run", "wshd.sock")

	user := "vcap"
	if spec.Privileged {
		user = "root"
	}

	wsh := exec.Command(wshPath, "--socket", sockPath, "--user", user, "/bin/bash")

	wsh.Stdin = bytes.NewBufferString(spec.Script)

	return c.jobTracker.Spawn(wsh)
}

func (c *LinuxContainer) Stream(jobID uint32) (<-chan backend.JobStream, error) {
	log.Println(c.id, "streaming job", jobID)
	return c.jobTracker.Stream(jobID)
}

func (c *LinuxContainer) Link(jobID uint32) (backend.JobResult, error) {
	log.Println(c.id, "linking to job", jobID)

	exitStatus, stdout, stderr, err := c.jobTracker.Link(jobID)
	if err != nil {
		return backend.JobResult{}, err
	}

	return backend.JobResult{
		ExitStatus: exitStatus,
		Stdout:     stdout,
		Stderr:     stderr,
	}, nil
}

func (c *LinuxContainer) NetIn(hostPort uint32, containerPort uint32) (uint32, uint32, error) {
	if hostPort == 0 {
		randomPort, err := c.portPool.Acquire()
		if err != nil {
			panic("x")
		}

		hostPort = randomPort
	}

	if containerPort == 0 {
		containerPort = hostPort
	}

	log.Println(
		c.id,
		"mapping host port",
		hostPort,
		"to container port",
		containerPort,
	)

	net := exec.Command(path.Join(c.path, "net.sh"), "in")

	net.Env = []string{
		fmt.Sprintf("HOST_PORT=%d", hostPort),
		fmt.Sprintf("CONTAINER_PORT=%d", containerPort),
	}

	return hostPort, containerPort, c.runner.Run(net)
}

func (c *LinuxContainer) NetOut(network string, port uint32) error {
	net := exec.Command(path.Join(c.path, "net.sh"), "out")

	if port != 0 {
		log.Println(
			c.id,
			"permitting traffic to",
			network,
			"with port",
			port,
		)

		net.Env = []string{
			"NETWORK=" + network,
			fmt.Sprintf("PORT=%d", port),
		}
	} else {
		if network == "" {
			return fmt.Errorf("network and/or port must be provided")
		}

		log.Println(c.id, "permitting traffic to", network)

		net.Env = []string{
			"NETWORK=" + network,
			"PORT=",
		}
	}

	return c.runner.Run(net)
}

func (c *LinuxContainer) rsync(src, dst string) error {
	wshPath := path.Join(c.path, "bin", "wsh")
	sockPath := path.Join(c.path, "run", "wshd.sock")

	rsync := exec.Command(
		"rsync",
		"-e", wshPath+" --socket "+sockPath+" --rsh",
		"-r",
		"-p",
		"--links",
		src,
		dst,
	)

	return c.runner.Run(rsync)
}

func (c *LinuxContainer) startOomNotifier() error {
	c.oomLock.Lock()
	defer c.oomLock.Unlock()

	if c.oomNotifier != nil {
		return nil
	}

	oomPath := path.Join(c.path, "bin", "oom")

	c.oomNotifier = exec.Command(oomPath, c.cgroupsManager.SubsystemPath("memory"))

	err := c.runner.Start(c.oomNotifier)
	if err != nil {
		return err
	}

	go c.watchForOom(c.oomNotifier)

	return nil
}

func (c *LinuxContainer) stopOomNotifier() {
	c.oomLock.RLock()
	defer c.oomLock.RUnlock()

	if c.oomNotifier != nil {
		c.runner.Kill(c.oomNotifier)
	}
}

func (c *LinuxContainer) watchForOom(oom *exec.Cmd) {
	err := c.runner.Wait(oom)
	if err == nil {
		log.Println(c.id, "out of memory")
		c.Stop(false)
	} else {
		log.Println(c.id, "oom failed:", err)
	}

	// TODO: handle case where oom notifier itself failed? kill container?
}

func parseMemoryStat(contents string) (stat backend.ContainerMemoryStat) {
	scanner := bufio.NewScanner(strings.NewReader(contents))

	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		field := scanner.Text()

		if !scanner.Scan() {
			break
		}

		value, err := strconv.ParseUint(scanner.Text(), 10, 0)
		if err != nil {
			continue
		}

		switch field {
		case "cache":
			stat.Cache = value
		case "rss":
			stat.Rss = value
		case "mapped_file":
			stat.MappedFile = value
		case "pgpgin":
			stat.Pgpgin = value
		case "pgpgout":
			stat.Pgpgout = value
		case "swap":
			stat.Swap = value
		case "pgfault":
			stat.Pgfault = value
		case "pgmajfault":
			stat.Pgmajfault = value
		case "inactive_anon":
			stat.InactiveAnon = value
		case "active_anon":
			stat.ActiveAnon = value
		case "inactive_file":
			stat.InactiveFile = value
		case "active_file":
			stat.ActiveFile = value
		case "unevictable":
			stat.Unevictable = value
		case "hierarchical_memory_limit":
			stat.HierarchicalMemoryLimit = value
		case "hierarchical_memsw_limit":
			stat.HierarchicalMemswLimit = value
		case "total_cache":
			stat.TotalCache = value
		case "total_rss":
			stat.TotalRss = value
		case "total_mapped_file":
			stat.TotalMappedFile = value
		case "total_pgpgin":
			stat.TotalPgpgin = value
		case "total_pgpgout":
			stat.TotalPgpgout = value
		case "total_swap":
			stat.TotalSwap = value
		case "total_pgfault":
			stat.TotalPgfault = value
		case "total_pgmajfault":
			stat.TotalPgmajfault = value
		case "total_inactive_anon":
			stat.TotalInactiveAnon = value
		case "total_active_anon":
			stat.TotalActiveAnon = value
		case "total_inactive_file":
			stat.TotalInactiveFile = value
		case "total_active_file":
			stat.TotalActiveFile = value
		case "total_unevictable":
			stat.TotalUnevictable = value
		}
	}

	return
}

func parseCPUStat(usage, statContents string) (stat backend.ContainerCPUStat) {
	cpuUsage, err := strconv.ParseUint(strings.Trim(usage, "\n"), 10, 0)
	if err != nil {
		return
	}

	stat.Usage = cpuUsage

	scanner := bufio.NewScanner(strings.NewReader(statContents))

	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		field := scanner.Text()

		if !scanner.Scan() {
			break
		}

		value, err := strconv.ParseUint(scanner.Text(), 10, 0)
		if err != nil {
			continue
		}

		switch field {
		case "user":
			stat.User = value
		case "system":
			stat.System = value
		}
	}

	return
}