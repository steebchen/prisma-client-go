package cmd

import (
	"errors"
	"fmt"
	"os/exec"
)

func Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		var exit *exec.ExitError
		if errors.Is(err, exit) {
			if exit.Success() {
				return nil
			}
			return fmt.Errorf("%s %s exited with status code %d and output %s: %w", name, args, exit.ExitCode(), string(out), err)
		}

		return fmt.Errorf("command %s %s failed: %w", name, args, err)
	}

	return nil
}
