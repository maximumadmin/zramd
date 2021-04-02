package system

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// IsRoot will check if current process was started by the init system (e.g.
// systemd) from which we expect to handle this process' capabilities, otherwise
// check if the current process is running as root.
func IsRoot() bool {
	return os.Getppid() == 1 || os.Geteuid() == 0
}

func cpuInfo() []byte {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		panic(err)
	}
	return data
}

// IsVM detects if we are currently running inside a VM, see also
// https://man.archlinux.org/man/systemd-detect-virt.1.en.
func IsVM() bool {
	// Try to run systemd-detect-virt which is more accurate but is not present on
	// all systems.
	cmd := exec.Command("systemd-detect-virt", "--vm")
	out, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(out)) != "none"
	}
	// If error happened because systemd-detect-virt is not available on the
	// system, try to use cpuinfo (less accurate but available everywhere).
	if errors.Is(err, exec.ErrNotFound) {
		info := cpuInfo()
		pattern := "(?m)^flags\\s*\\:.*\\s+hypervisor(?:\\s+.*)?$"
		match, _ := regexp.Match(pattern, info)
		return match
	}
	panic(err)
}
