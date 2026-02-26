package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/steebchen/prisma-client-go/cli"
	"github.com/steebchen/prisma-client-go/logger"
)

func main() {
	if len(os.Args) > 1 {
		args := os.Args[1:]
		logger.Debug.Printf("invoking command %+v", args)

		switch args[0] {
		case "prefetch":
			 just run prisma -v to trigger the download
			if err := cli.Run([]string{"-v"}, true); err != nil {
				panic(err)
			}
			os.Exit(0)
			return
		case "init":
			 override default init flags
			args = append(args, "--generator-provider", "go run github.com/steebchen/prisma-client-go")
			if err := cli.Run(args, true); err != nil {
				panic(err)
			}
			os.Exit(0)
			return
		}

		 prisma CLI
		if err := cli.Run(args, true); err != nil {
			panic(err)
		}

		return
	}

	 running the prisma generator

	logger.Debug.Printf("invoking prisma")

 if this wasn't actually invoked by the prisma generator, print a warning and exit
	if os.Getenv("PRISMA_GENERATOR_INVOCATION") == "" {
		logger.Info.Printf("This command is only meant to be invoked internally. Please run the following instead:")
		logger.Info.Printf("`go run github.com/steebchen/prisma-client-go <command>`")
		logger.Info.Printf("e.g.")
		logger.Info.Printf("`go run github.com/steebchen/prisma-client-go generate`")
		os.Exit(1)
	}

	 exit when signal triggers
	c := make(chan os.Signal, 1)
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
