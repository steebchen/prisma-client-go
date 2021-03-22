package transaction

import (
	"encoding/json"
)

type Result struct{}

func (r *Result) Get(c <-chan []byte, v interface{}) error {
	data, ok := <-c
	if !ok {
		return nil
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}
