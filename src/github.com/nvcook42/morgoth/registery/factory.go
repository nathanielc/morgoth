package registery

import (
	"github.com/nvcook42/morgoth/config/types"
)

// A Factory provides methods to create an empty configuration object
// and to get an instance from a populated configuration
type Factory interface {
	//Return a new zeroed configuration.
	NewConf() types.Configuration
	//Return a new instance based on the given configuration.
	GetInstance(types.Configuration) (interface{}, error)
}
