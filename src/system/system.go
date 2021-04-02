package system

import (
	"os"
	"os/exec"
	"strings"
)

// IsRoot will check if current process was started by the init system (e.g.
// systemd) from which we expect to handle this process' capabilities, otherwise
// check if the current process is running as root.
func IsRoot() bool {
	return os.Getppid() == 1 || os.Geteuid() == 0
}

// IsVM detects if we are currently running inside a VM, if systemd-detect-virt
// is missing (i.e. on non systemd systems), the result will be false, see also
// https://man.archlinux.org/man/systemd-detect-virt.1.en.
func IsVM() bool {
	cmd := exec.Command("systemd-detect-virt", "--vm")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) != "none"
}
