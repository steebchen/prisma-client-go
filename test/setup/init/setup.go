package main

import (
	"log"

	"github.com/prisma/prisma-client-go/test"
)

func main() {
	log.Printf("setting up tests")
	teardown()
	setup()

	// teardown()
}

func setup() {
	for _, db := range test.Databases {
		db.Setup()
	}

	log.Printf("setup done")
}

func teardown() {
	for _, db := range test.Databases {
		db.Teardown()
	}
	log.Printf("teardown done")
}
