package config

import (
	"github.com/nvcook42/morgoth/validator"
)

func (self Config) Validate() error {
	return validator.ValidateAll(self)
}

func (self DataEngine) Validate() error {
	return validator.ValidateOne(self)
}

func (self Metrics) Validate() error {
	return validator.ValidateAll(self)
}

func (self Fittings) Validate() error {
	return validator.ValidateAll(self)
}
