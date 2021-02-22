package util

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
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
