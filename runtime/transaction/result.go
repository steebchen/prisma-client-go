package transaction

import (
	"encoding/json"
	"fmt"

	"github.com/steebchen/prisma-client-go/engine"
	"github.com/steebchen/prisma-client-go/logger"
)

type Result struct {
	cache []byte
}

func (r *Result) Get(c <-chan []byte, v interface{}) error {
	var res []byte
	if r.cache != nil {
		res = r.cache
	} else {
		data, ok := <-c
		if !ok {
			return fmt.Errorf("result not fetched")
		}
		data, err := engine.TransformResponse(data)
		if err != nil {
			return fmt.Errorf("could not transform response: %w", err)
		}
		res = data
		r.cache = data
	}
	logger.Debug.Printf("tx result: %s", res)
	if err := json.Unmarshal(res, &v); err != nil {
		return err
	}
	return nil
}
