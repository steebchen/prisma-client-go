package cmd

import (
	"fmt"
	"os/exec"
)

func Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		exit, ok := err.(*exec.ExitError)
		if !ok {
			return fmt.Errorf("command %s %s failed: %w", name, args, err)
		}

		if !exit.Success() {
			return fmt.Errorf("%s %s exited with status code %d and output %s: %w", name, args, exit.ExitCode(), string(out), err)
		}
	}

	return nil
}
