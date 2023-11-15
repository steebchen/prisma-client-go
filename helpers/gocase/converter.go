package gocase

// Converter represents a conversion object.
// If you need a conversion other than the default,
// you need create a new conversion object.
type Converter struct {
	initialisms []initialism
}

var defaultConverter = newDefault()

func newDefault() *Converter {
	i, _ := createInitialisms(DefaultInitialisms...)
	return &Converter{
		initialisms: i,
	}
}

// New creates a new Converter.
func New(opts ...Option) (*Converter, error) {
	c := newDefault()
	for _, opt := range opts {
		if err := opt.apply(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// Option is an option that configures Converter.
type Option interface {
	apply(*Converter) error
}

type funcOption struct {
	f func(*Converter) error
}

func (f *funcOption) apply(o *Converter) error {
	return f.f(o)
}

func newFuncOption(f func(*Converter) error) *funcOption {
	return &funcOption{
		f: f,
	}
}

// WithInitialisms is an option to set initialisms.
// If you want to add to the default initialisms, use `DefaultInitialisms`.
func WithInitialisms(initialisms ...string) Option {
	return newFuncOption(func(c *Converter) error {
		i, err := createInitialisms(initialisms...)
		if err != nil {
			return err
		}
		c.initialisms = i
		return nil
	})
}
