package config

import (
	"errors"
	"github.com/nvcook42/morgoth/fittings/rest"
	"github.com/nvcook42/morgoth/validator"
)

type FittingType string

const (
	REST     FittingType = "rest"
	Graphite FittingType = "graphite"
)

func (self FittingType) Validate() error {
	switch self {
	case
		REST,
		Graphite:
		return nil
	}
	return errors.New("Invalid FittingType")
}

// The Fittings subsection of the config
type Fitting struct {
	Type FittingType
	REST *rest.RESTConf `yaml:"rest,omitempty"`
	//	Graphite *graphite.GraphiteConf `yaml:"graphite,omitempty"`
}

func (self *Fitting) Default() {
	switch self.Type {
	case REST:
		//self.REST.Default()
	}
}

func (self Fitting) Validate() error {
	return validator.ValidateOne(self)
}
