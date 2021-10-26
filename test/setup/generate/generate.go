// Package generate runs prisma generate in parallel
package main

//go:generate go run .

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		if !strings.Contains(path, "migrations") && info.Name() == "schema.prisma" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// prefetch the binaries first
	// generating in this folder makes sure binaries are unpacked
	cmd := exec.Command("go", "run", "github.com/prisma/prisma-client-go", "generate")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	// manually unpack binary
	log.Printf("unpacking binaries...")
	cmd = exec.Command("go", "run", "./db")
	cmd.Env = append(os.Environ(), "PRISMA_CLIENT_GO_LOG=info")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	maxGoroutines := 5
	throttle := make(chan struct{}, maxGoroutines)

	for _, file := range files {
		wg.Add(1)
		throttle <- struct{}{}
		go func(file string) {
			defer wg.Done()

			genCmd := exec.Command("go", "run", "github.com/prisma/prisma-client-go", "generate")
			genCmd.Dir = filepath.Dir(file)
			genCmd.Stderr = os.Stderr
			genCmd.Stdout = os.Stdout
			if err := genCmd.Run(); err != nil {
				log.Fatal(err)
			}

			log.Printf("%s done", file)

			<-throttle
		}(file)
	}

	wg.Wait()
}
