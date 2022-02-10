package initializer

import (
	"github.com/hashicorp/go-multierror"

	"github.com/mrlyc/cmdr/cmdr"
)

type Chaining struct {
	initializers []cmdr.Initializer
}

func (c *Chaining) Add(target ...interface{}) *Chaining {
	for _, t := range target {
		initializer, ok := t.(cmdr.Initializer)
		if !ok {
			continue
		}

		c.initializers = append(c.initializers, initializer)
	}

	return c
}

func (c Chaining) GetInitializers() []cmdr.Initializer {
	return c.initializers
}

func (c *Chaining) Init() error {
	var errs error
	for _, initializer := range c.initializers {
		err := initializer.Init()
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

func NewChaining(initializers ...cmdr.Initializer) *Chaining {
	return &Chaining{
		initializers: initializers,
	}
}
