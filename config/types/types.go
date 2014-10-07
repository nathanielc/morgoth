// Common interfaces needed for configuration
package types

import (
	"github.com/nvcook42/morgoth/defaults"
)

type Validator interface {
	Validate() error
}

type Configuration interface {
	defaults.Defaulter
	Validator
}
