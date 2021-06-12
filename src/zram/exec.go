package zram

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func execute(command string, arg ...string) error {
	cmd := exec.Command(command, arg...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if len(msg) == 0 {
			msg = fmt.Sprintf(
				"failed to execute \"%s\"",
				strings.Join(append([]string{command}, arg...), " "),
			)
		}
		return errors.New(msg)
	}
	return nil
}
