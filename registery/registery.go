// The registry package provides a simple plugin registery framework.
// The Registery struct provides method to register and retrieve
// factories that can later be used to dynamically create instances
// of any type.
package registery

import (
	"errors"
	"fmt"
	"github.com/golang/glog"
)

// A Registery for mapping names to factories
type Registery struct {
	factories map[string]Factory
}

// Create a new Registery object
func New() *Registery {
	r := new(Registery)
	r.factories = make(map[string]Factory)
	return r
}

// Register a Factory by name
func (self *Registery) RegisterFactory(name string, factory Factory) error {
	glog.V(2).Infof("Registering Factory %s", name)
	if _, ok := self.factories[name]; ok {
		return errors.New(fmt.Sprintf("Factory of name %s already registered"))
	}
	self.factories[name] = factory
	return nil
}

// Get a registered Factory by name
func (self *Registery) GetFactory(name string) (Factory, error) {
	factory, ok := self.factories[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Unknown Factory %s", name))
	}
	return factory, nil
}
