package networking_test

import (
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cloudfoundry-incubator/garden/api"
	"github.com/cloudfoundry/gunk/localip"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("IP settings", func() {
	var (
		container          api.Container
		containerNetwork   string
		containerInterface string
		hostInterface      string
	)

	JustBeforeEach(func() {
		client = startGarden()

		var err error

		container, err = client.Create(api.ContainerSpec{Network: containerNetwork})
		Ω(err).ShouldNot(HaveOccurred())

		containerInterface = "w" + strconv.Itoa(GinkgoParallelNode()) + container.Handle() + "-1"
		hostInterface = "w" + strconv.Itoa(GinkgoParallelNode()) + container.Handle() + "-0"
	})

	AfterEach(func() {
		err := client.Destroy(container.Handle())
		Ω(err).ShouldNot(HaveOccurred())
	})

	Context("when the Network parameter is a subnet address", func() {
		BeforeEach(func() {
			containerNetwork = "10.3.0.0/24"
		})

		Describe("container's network interface", func() {
			It("has the correct IP address", func() {
				stdout := gbytes.NewBuffer()
				stderr := gbytes.NewBuffer()

				process, err := container.Run(api.ProcessSpec{
					Path: "/sbin/ifconfig",
					Args: []string{containerInterface},
				}, api.ProcessIO{
					Stdout: stdout,
					Stderr: stderr,
				})
				Ω(err).ShouldNot(HaveOccurred())
				rc, err := process.Wait()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(rc).Should(Equal(0))

				Ω(stdout.Contents()).Should(ContainSubstring(" inet addr:10.3.0.1 "))
			})
		})

		Describe("hosts's network interface for a container", func() {
			It("has the correct IP address", func() {

				out, err := exec.Command("/sbin/ifconfig", hostInterface).Output()
				Ω(err).ShouldNot(HaveOccurred())

				Ω(out).Should(ContainSubstring(" inet addr:10.3.0.254 "))
			})
		})
	})

	Context("when the Network parameter is not a subnet address", func() {
		BeforeEach(func() {
			containerNetwork = "10.3.0.2/24"
		})

		Describe("container's network interface", func() {
			It("has the specified IP address", func() {
				stdout := gbytes.NewBuffer()
				stderr := gbytes.NewBuffer()

				process, err := container.Run(api.ProcessSpec{
					Path: "/sbin/ifconfig",
					Args: []string{containerIfName(container)},
				}, api.ProcessIO{
					Stdout: stdout,
					Stderr: stderr,
				})
				Ω(err).ShouldNot(HaveOccurred())
				rc, err := process.Wait()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(rc).Should(Equal(0))

				Ω(stdout.Contents()).Should(ContainSubstring(" inet addr:10.3.0.2 "))
			})
		})

		Describe("hosts's network interface for a container", func() {
			It("has the correct IP address", func() {

				out, err := exec.Command("/sbin/ifconfig", hostInterface).Output()
				Ω(err).ShouldNot(HaveOccurred())

				Ω(out).Should(ContainSubstring(" inet addr:10.3.0.254 "))
			})
		})
	})

	Describe("the container's network", func() {
		It("is reachable from the host", func() {
			info, ierr := container.Info()
			Ω(ierr).ShouldNot(HaveOccurred())

			out, err := exec.Command("/bin/ping", "-c 2", info.ContainerIP).Output()
			Ω(err).ShouldNot(HaveOccurred())

			Ω(out).Should(ContainSubstring(" 0% packet loss"))
		})
	})

	Describe("host's network", func() {
		It("is reachable from inside the container", func() {
			info, ierr := container.Info()
			Ω(ierr).ShouldNot(HaveOccurred())

			stdout := gbytes.NewBuffer()
			stderr := gbytes.NewBuffer()

			listener, err := net.Listen("tcp", fmt.Sprintf("%s:0", info.HostIP))
			Ω(err).ShouldNot(HaveOccurred())
			defer listener.Close()

			mux := http.NewServeMux()
			mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hello")
			})

			go (&http.Server{Handler: mux}).Serve(listener)

			process, err := container.Run(api.ProcessSpec{
				Path: "sh",
				Args: []string{"-c", fmt.Sprintf("(echo 'GET /test HTTP/1.1'; echo 'Host: foo.com'; echo) | nc %s %s", info.HostIP, strings.Split(listener.Addr().String(), ":")[1])},
			}, api.ProcessIO{
				Stdout: stdout,
				Stderr: stderr,
			})
			Ω(err).ShouldNot(HaveOccurred())

			rc, err := process.Wait()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(rc).Should(Equal(0))

			Ω(stdout.Contents()).Should(ContainSubstring("Hello"))
		})
	})

	Describe("the container's external ip", func() {
		It("is the external IP of its host", func() {
			info, err := container.Info()
			Ω(err).ShouldNot(HaveOccurred())

			localIP, err := localip.LocalIP()
			Ω(err).ShouldNot(HaveOccurred())

			Ω(localIP).Should(Equal(info.ExternalIP))
		})
	})
})