package main

import (
	"log"
	"os"
	"time"

	"github.com/prisma/prisma-client-go/test"
)

func main() {
	action := ""
	if len(os.Args) > 1 {
		action = os.Args[1]
	}

	switch action {
	case "setup":
		teardown()
		setup()
	case "teardown":
		teardown()
	default:
		log.Fatalf("no such action %s, only 'setup' or 'teardown' are accepted", action)
	}
}

func setup() {
	for _, db := range test.Databases {
		db.Setup()
	}

	time.Sleep(15 * time.Second)

	log.Printf("setup done")
}

func teardown() {
	for _, db := range test.Databases {
		db.Teardown()
	}
	log.Printf("teardown done")
}
