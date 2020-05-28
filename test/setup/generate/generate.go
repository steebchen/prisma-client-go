// Package generate runs prisma generate in parallel
package main

//go:generate go run .

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func main() {
	log.Printf("generating clients")

	generate()
}

func generate() {
	var wg sync.WaitGroup

	var files []string
	err := filepath.Walk("../..", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "schema.prisma" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// prefetch the binaries first
	cmd := exec.Command("go", "run", "github.com/prisma/prisma-client-go", "prefetch")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()

			cmd := exec.Command("go", "run", "github.com/prisma/prisma-client-go", "generate")
			cmd.Dir = filepath.Dir(file)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}

			log.Printf("%s done", file)
		}(file)
	}

	wg.Wait()
}
