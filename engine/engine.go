package engine

import (
	"context"
)

type Engine interface {
	Connect() error
	Disconnect() error
	Do(ctx context.Context, payload interface{}, into interface{}) error
	Batch(ctx context.Context, payload interface{}, into interface{}) error
	Name() string
}
