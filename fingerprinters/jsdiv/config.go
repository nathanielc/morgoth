package jsdiv

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	"github.com/nathanielc/morgoth/config"
)

type JSDivConf struct {
	Min    float64
	Max    float64
	NBins  int     `validate:"min=1"`
	PValue float64 `validate:"nonzero,min=0,max=1"`
}

func (self *JSDivConf) Validate() error {
	return validator.Validate(self)
}

func (self *JSDivConf) Default() {
	config.PerformDefault(self)
}
