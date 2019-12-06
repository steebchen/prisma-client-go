package main

import (
	"log"
	"os"

	"github.com/prisma/photongo/generate"
	"github.com/prisma/photongo/logger"
)

func main() {
	log.Printf("args: %+v", os.Args)

	cmd := ""

	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	if cmd == "generate" {
		logger.L.Printf("invoking generate")
		// prisma CLI
		err := generate.Run()
		if err != nil {
			panic(err)
		}
	} else {
		logger.L.Printf("invoking prisma")
		// invoke the prisma generator
		invokePrisma()
	}
}
