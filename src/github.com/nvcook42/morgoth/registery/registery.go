package registery

import (
	"errors"
	"fmt"
)

type Registery struct {
	factories map[string]Factory
}

func New() *Registery {
	r := new(Registery)
	r.factories = make(map[string]Factory)
	return r
}

func (self *Registery) RegisterFactory(name string, factory Factory) error {
	if _, ok := self.factories[name]; ok {
		return errors.New(fmt.Sprintf("Factory of name %s already registered"))
	}
	self.factories[name] = factory
	return nil
}

func (self *Registery) GetFactory(name string) (Factory, error) {
	factory, ok := self.factories[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Unknown Factory %s", name))
	}
	return factory, nil
}
