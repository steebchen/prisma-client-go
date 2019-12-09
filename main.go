package main

import (
	"os"

	"github.com/prisma/photongo/cli"
	"github.com/prisma/photongo/logger"
)

func main() {
	if len(os.Args) > 1 {
		args := os.Args[1:]
		logger.L.Printf("invoking command %+v", args)
		// prisma CLI
		err := cli.Run(args)
		if err != nil {
			panic(err)
		}
	} else {
		logger.L.Printf("invoking prisma")
		// invoke the prisma generator
		invokePrisma()
	}
}
