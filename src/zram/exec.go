package zram

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

func execute(command string, arg ...string) error {
	cmd := exec.Command(command, arg...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return errors.New(strings.TrimSpace(stderr.String()))
	}
	return nil
}
