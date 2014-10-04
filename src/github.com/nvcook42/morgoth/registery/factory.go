package registery

import (
	"github.com/nvcook42/morgoth/config/types"
)
type Factory interface {
	NewConf() types.Configuration
	GetInstance(types.Configuration) (interface{}, error)
}
