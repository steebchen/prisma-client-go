package hooks

import (
	"context"
	"fmt"
	"os/exec"
	"testing"

	"github.com/prisma/photongo/cli"
)

type gqlResponse struct {
	Data   interface{} `json:"data"`
	Errors []gqlError  `json:"errors"`
}

type gqlError struct {
	Message string `json:"error"`
}

type TestClient interface {
	Connect() error
	Disconnect() error
}

func Start(t *testing.T, client TestClient, before string, do func(context.Context, string, interface{}) error) {
	setup(t)

	if err := client.Connect(); err != nil {
		t.Fatalf("could not connect: %s", err)
		return
	}

	ctx := context.Background()

	if before != "" {
		var response gqlResponse
		err := do(ctx, before, &response)
		if err != nil {
			t.Fatalf("could not send mock query %s", err)
		}
		if response.Errors != nil {
			t.Fatalf("mock query has errors %+v", response)
		}
	}
}

func End(t *testing.T, client TestClient) {
	err := client.Disconnect()
	if err != nil {
		t.Fatalf("could not disconnect: %s", err)
	}
}

func setup(t *testing.T) {
	if err := cmd("rm", "-rf", "dev.sqlite"); err != nil {
		t.Fatal(err)
	}

	if err := cmd("rm", "-rf", "migrations"); err != nil {
		t.Fatal(err)
	}

	if err := cli.Run([]string{"lift", "save", "--create-db", "--name", "init"}); err != nil {
		t.Fatalf("could not run lift save %s", err)
	}

	if err := cli.Run([]string{"lift", "up"}); err != nil {
		t.Fatalf("could not run lift save %s", err)
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
