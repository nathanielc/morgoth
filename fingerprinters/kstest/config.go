package kstest

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	"github.com/nathanielc/morgoth/config"
)

type KSTestConf struct {
	Confidence uint `validate:"min=1,max=5 default:"4"`
}

func (self *KSTestConf) Validate() error {
	return validator.Validate(self)
}

func (self *KSTestConf) Default() {
	config.PerformDefault(self)
}
