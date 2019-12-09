package hooks

import (
	"fmt"
	"log"
	"os/exec"
	"testing"
)

const pkg = "github.com/prisma/photongo"

func Run(t *testing.T) {
	if err := cmd("rm", "-rf", "dev.sqlite"); err != nil {
		log.Fatal(err)
	}
	if err := cmd("rm", "-rf", "migrations"); err != nil {
		log.Fatal(err)
	}

	if err := cmd("go", "run", pkg, "lift", "save", "--create-db", "--name", "init"); err != nil {
		t.Fatalf("could not run lift save %s", err)
	}
	if err := cmd("go", "run", pkg, "lift", "up"); err != nil {
		t.Fatalf("could not run lift up %s", err)
	}
}

func cmd(name string, args ...string) error {
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
