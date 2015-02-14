package graphite

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	config "github.com/nvcook42/morgoth/config/types"
)

type GraphiteConf struct {
	Port uint `validate:"min=1,max=65535" default:"2003"`
}

func (self *GraphiteConf) Validate() error {
	return validator.Validate(self)
}

//Sets any invalid fields to their default value
func (self *GraphiteConf) Default() {
	config.PerformDefault(self)
}
