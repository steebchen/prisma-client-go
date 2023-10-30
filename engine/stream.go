package engine

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

type Messsage struct {
	IsPanic bool   `json:"is_panic"`
	Message string `json:"message"`
}

func checkStderr(cmd *exec.Cmd) error {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("get stderr pipe: %w", err)
	}

	go func() {
		scanner := bufio.NewScanner(stderr)
		const maxCapacity int = 65536
		buf := make([]byte, maxCapacity)
		scanner.Buffer(buf, maxCapacity)

		// optionally, resize scanner's capacity for lines over 64K, see next example
		for scanner.Scan() {
			contents := scanner.Bytes()

			var message Messsage
			if err := json.Unmarshal(contents, &message); err != nil {
				log.Printf("failed to unmarshal message: %s", err.Error())
			}

			if message.Message != "" {
				log.Println(message.Message)
				continue
			}

			log.Println(string(contents))
		}
	}()

	return nil
}
