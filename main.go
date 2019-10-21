package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/prisma/photongo/generator"
	"github.com/prisma/photongo/jsonrpc"
)

func reply(w io.Writer, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "could not marshal data")
	}
	_, err = w.Write(b)
	if err != nil {
		return errors.Wrap(err, "could not write data")
	}
	_, err = w.Write([]byte("\n"))
	if err != nil {
		return errors.Wrap(err, "could not write data")
	}
	return nil
}

func main() {
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
			response = NewManifest()

		case "generate":
			response = nil // success

			var params generator.Root
			err := mapstructure.Decode(input.Params, &params)
			if err != nil {
				log.Fatalf("could not assert params into generator.Options type %s", err)
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
