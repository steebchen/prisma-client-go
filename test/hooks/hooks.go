package hooks

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"testing"

	"github.com/prisma/prisma-client-go/cli"
	"github.com/prisma/prisma-client-go/engine"
	"github.com/prisma/prisma-client-go/logger"
)

type Engine interface {
	Connect() error
	Disconnect() error
	Do(context.Context, string, interface{}) error
}

func Start(t *testing.T, e *engine.Engine, before []string) {
	setup(t)

	if err := e.Connect(); err != nil {
		t.Fatalf("could not connect: %s", err)
		return
	}

	ctx := context.Background()

	for _, b := range before {
		var response engine.GQLResponse
		err := e.Do(ctx, b, &response)
		if err != nil {
			End(t, e)
			t.Fatalf("could not send mock query %s", err)
		}
		if response.Errors != nil {
			End(t, e)
			t.Fatalf("mock query has errors %+v", response)
		}
	}

	log.Printf("")
	log.Printf("---")
	log.Printf("")
}

func End(t *testing.T, e Engine) {
	defer cleanup(t)

	err := e.Disconnect()
	if err != nil {
		t.Fatalf("could not disconnect: %s", err)
	}
}

func setup(t *testing.T) {
	cleanup(t)

	if err := cli.Run([]string{"migrate", "save", "--experimental", "--create-db", "--name", "init"}, logger.Enabled); err != nil {
		t.Fatalf("could not run migrate save --experimental %s", err)
	}

	if err := cli.Run([]string{"migrate", "up", "--experimental"}, logger.Enabled); err != nil {
		t.Fatalf("could not run migrate save --experimental %s", err)
	}
}

func cleanup(t *testing.T) {
	if err := Cmd("rm", "-rf", "dev.sqlite"); err != nil {
		t.Fatal(err)
	}

	if err := Cmd("rm", "-rf", "migrations"); err != nil {
		t.Fatal(err)
	}
}

func Cmd(name string, args ...string) error {
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
