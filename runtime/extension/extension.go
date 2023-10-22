package extension

type Extension[T client] struct {
	Client *T
}

func (r *Extension[T]) Extend(extension ...Action) *T {
	return r.Client
}

type Result struct {
}

type Run = func() (*Result, error)

type Action struct {
}

type Args map[string]interface{}

type client interface {
	// TODO
	//Engine() *engine.Engine
}
