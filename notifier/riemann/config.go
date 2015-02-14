package riemann

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	config "github.com/nvcook42/morgoth/config/types"
)

type RiemannConf struct {
	Host string `validate:"min=1" default:"localhost"`
	Port uint   `validate:"min=1,max=65535" default:"5555"`
}

func (self *RiemannConf) Validate() error {
	return validator.Validate(self)
}

//Sets any invalid fields to their default value
func (self *RiemannConf) Default() {
	config.PerformDefault(self)
}
