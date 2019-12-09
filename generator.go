package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mitchellh/mapstructure"

	"github.com/prisma/photongo/generator"
	"github.com/prisma/photongo/jsonrpc"
)

func reply(w io.Writer, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("could not marshal data %s", err)
	}

	b = append(b, byte('\n'))

	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("could not write data %s", err)
	}

	return nil
}

func invokePrisma() {
	// make sure to exit when signal triggers
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		os.Exit(1)
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		content := scanner.Bytes()

		var input jsonrpc.Request

		err := json.Unmarshal(content, &input)
		if err != nil {
			log.Fatalf("could not open stdin %s", err)
		}

		var response interface{}

		switch input.Method {
		case "getManifest":
			response = jsonrpc.ManifestResponse{
				Manifest: jsonrpc.Manifest{
					DefaultOutput:      "./photon/photon_gen.go",
					PrettyName:         "Photon Go",
					Denylist:           []string{},
					RequiresGenerators: []string{},
					RequiresEngines:    []string{}, // Photon Go handles downloading the engines
				},
			}

		case "generate":
			response = nil // success

			var params generator.Root

			err := mapstructure.Decode(input.Params, &params)
			if err != nil {
				log.Fatalf("could not assert params into generator.Root type %s", err)
			}

			err = generator.Run(params)
			if err != nil {
				log.Fatalf("could not generate code. %s", err)
			}
		}

		err = reply(os.Stderr, jsonrpc.NewResponse(input.ID, response))

		if err != nil {
			log.Fatalf("could not open stdin %s", err)
		}
	}
}
