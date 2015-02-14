package kstest

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	config "github.com/nvcook42/morgoth/config/types"
)

type KSTestConf struct {
	Confidence      uint `yaml:"confidence" validate:"min=1,max=5" default:"1"`
	NormalCount     uint `yaml:"normal_count" validate:"nonzero" default:"3"`
	MaxFingerprints uint `yaml:"max_fingerprints" validate:"nonzero" default:"20"`
}

func (self *KSTestConf) Default() {
	config.PerformDefault(self)
}

func (self *KSTestConf) Validate() error {
	return validator.Validate(self)
}
