package mock

import (
	"sync"
)

func New(expectations *[]Expectation) *Engine {
	return &Engine{
		expectations: expectations,
	}
}

type Engine struct {
	expectations *[]Expectation
	expMu        sync.Mutex
}

func (e *Engine) Name() string {
	return "mock"
}

func (e *Engine) Connect() error {
	panic("this is a mock client – you don't need to connect or disconnect this client")
}

func (e *Engine) Disconnect() error {
	panic("this is a mock client – you don't need to connect or disconnect this client")
}
