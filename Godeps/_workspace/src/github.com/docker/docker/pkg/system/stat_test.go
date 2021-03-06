package system

import (
	"syscall"
	"testing"
)

func TestFromStatT(t *testing.T) {
	file, _, _ := prepareFiles(t)

	stat := &syscall.Stat_t{}
	err := syscall.Lstat(file, stat)

	s, err := fromStatT(stat)
	if err != nil {
		t.Fatal(err)
	}

	if stat.Mode != s.Mode() {
		t.Fatal("got invalid mode")
	}
	if stat.Uid != s.Uid() {
		t.Fatal("got invalid uid")
	}
	if stat.Gid != s.Gid() {
		t.Fatal("got invalid gid")
	}
	if stat.Rdev != s.Rdev() {
		t.Fatal("got invalid rdev")
	}
	if stat.Mtim != s.Mtim() {
		t.Fatal("got invalid mtim")
	}
}
