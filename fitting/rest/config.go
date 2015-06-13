package rest

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	config "github.com/nathanielc/morgoth/config/types"
)

type RESTConf struct {
	Port uint `validate:"min=1,max=65535" default:"8000"`
}

func (self *RESTConf) Validate() error {
	return validator.Validate(self)
}

//Sets any invalid fields to their default value
func (self *RESTConf) Default() {
	config.PerformDefault(self)
}
