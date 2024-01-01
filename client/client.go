package client

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/steebchen/prisma-client-go/cli"
	"github.com/steebchen/prisma-client-go/logger"
)

func Process(args []string) error {

	if len(args) == 0 {
		return cli.Run([]string{"--help"}, true)
	}

	switch args[0] {
	case "prefetch":
		return cli.Run([]string{"-v"}, true)
	case "init":
		// override default init flags
		args = append(args, "--generator-provider", ".")
		return cli.Run(args, true)
	}

	// exit when signal triggers
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		os.Exit(1)
	}()

	if err := invokePrisma(); err != nil {
		log.Printf("error occurred when invoking prisma: %s", err)
		return err
	}

	logger.Debug.Printf("success")

	return cli.Run(args, true)
}
