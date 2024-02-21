package engine

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/steebchen/prisma-client-go/logger"
)

type Messsage struct {
	IsPanic bool   `json:"is_panic"`
	Message string `json:"message"`
}

func (e *QueryEngine) streamStderr(cmd *exec.Cmd, onError chan<- string) error {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("get stderr pipe: %w", err)
	}

	go func() {
	outer:
		for {
			select {
			case v := <-e.onEngineError:
				e.mu.Lock()
				e.lastEngineError = v
				e.mu.Unlock()
			case <-e.closed:
				logger.Debug.Printf("query engine closed")
				break outer
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		const maxCapacity int = 65536
		buf := make([]byte, maxCapacity)
		scanner.Buffer(buf, maxCapacity)

		for scanner.Scan() {
			contents := scanner.Bytes()
			var message Messsage
			if err := json.Unmarshal(contents, &message); err != nil {
				log.Printf("failed to unmarshal message: %s", err.Error())
				continue
			}

			if message.Message != "" {
				onError <- message.Message
				log.Println(message.Message)
				continue
			}

			log.Println(string(contents))
		}
	}()

	return nil
}
