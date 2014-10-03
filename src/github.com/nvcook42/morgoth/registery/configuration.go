package registery

import (
	"github.com/nvcook42/morgoth/defaults"
	"github.com/nvcook42/morgoth/validator"
)

type Configuration interface {
	defaults.Defaulter
	validator.Validator
}
