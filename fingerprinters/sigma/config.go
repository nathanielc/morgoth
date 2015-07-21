package sigma

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	"github.com/nathanielc/morgoth/config"
)

type SigmaConf struct {
	Deviations float64 `validiate:"nonzero"`
}

func (self *SigmaConf) Validate() error {
	return validator.Validate(self)
}

func (self *SigmaConf) Default() {
	config.PerformDefault(self)
}
