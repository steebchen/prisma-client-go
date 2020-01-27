package main

//go:generate go build .

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/prisma/prisma-client-go/cli"
	"github.com/prisma/prisma-client-go/logger"
)

func main() {
	if len(os.Args) > 1 {
		args := os.Args[1:]
		logger.Debug.Printf("invoking command %+v", args)

		if args[0] == "prefetch" {
			// just run prisma -v to trigger the download
			if err := cli.Run([]string{"-v"}, true); err != nil {
				panic(err)
			}
			os.Exit(0)
			return
		}

		// prisma CLI
		if err := cli.Run(args, true); err != nil {
			panic(err)
		}

		return
	}

	// running the prisma generator

	logger.Debug.Printf("invoking prisma")

	// exit when signal triggers
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		os.Exit(1)
	}()

	if err := invokePrisma(); err != nil {
		log.Printf("error occurred when invoking prisma: %s", err)
		os.Exit(1)
	}

	logger.Debug.Printf("success")
}
