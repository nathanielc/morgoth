package tukey

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	config "github.com/nathanielc/morgoth/config/types"
)

type TukeyConf struct {
	Threshold float64 `yaml:"threshold" validate:"min=0,nonzero" default:"3"`
}

func (self *TukeyConf) Default() {
	config.PerformDefault(self)
}

func (self *TukeyConf) Validate() error {
	return validator.Validate(self)
}
