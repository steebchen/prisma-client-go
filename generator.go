package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/prisma/prisma-client-go/generator"
	"github.com/prisma/prisma-client-go/jsonrpc"
	"github.com/prisma/prisma-client-go/logger"
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

func invokePrisma() error {
	reader := bufio.NewReader(os.Stdin)

	if logger.Enabled {
		dir, _ := os.Getwd()
		log.Printf("current working dir: %s", dir)
	}

	for {
		content, err := reader.ReadBytes('\n')
		if err == io.EOF {
			logger.Debug.Printf("warning: ignoring EOF error. stdin: `%s`", content)
			return nil
		}
		if err != nil {
			return fmt.Errorf("could not read bytes from stdin: %w", err)
		}

		var input jsonrpc.Request

		if err := json.Unmarshal(content, &input); err != nil {
			return fmt.Errorf("could not open stdin %s", err)
		}

		var response interface{}

		switch input.Method {
		case "getManifest":
			response = jsonrpc.ManifestResponse{
				Manifest: jsonrpc.Manifest{
					DefaultOutput: path.Join(".", "db"),
					PrettyName:    "Prisma Client Go",
				},
			}

		case "generate":
			response = nil // success

			var params generator.Root

			if err := json.Unmarshal(input.Params, &params); err != nil {
				return fmt.Errorf("could not unmarshal params into generator.Root type %s", err)
			}

			if err = generator.Run(&params); err != nil {
				return fmt.Errorf("could not generate code. %s", err)
			}
		default:
			return fmt.Errorf("no such method %s", input.Method)
		}

		if err := reply(os.Stderr, jsonrpc.NewResponse(input.ID, response)); err != nil {
			return fmt.Errorf("could not reply %s", err)
		}
	}
}
