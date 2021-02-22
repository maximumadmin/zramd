package util

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"zramd/src/uname"
)

// Run executes a command and returns it's stderr as error if it failed.
func Run(command string, arg ...string) error {
	cmd := exec.Command(command, arg...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return errors.New(strings.TrimSpace(stderr.String()))
	}
	return nil
}

// IsRoot checks if program is running as root.
func IsRoot() bool {
	return os.Geteuid() == 0
}

// IsZstdSupported checks if current kernel supports zstd compressed zram.
func IsZstdSupported() bool {
	major, minor := uname.Uname().KernelVersion()
	return (major == 4 && minor >= 19) || major > 4
}
