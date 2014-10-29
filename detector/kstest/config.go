package kstest

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/defaults"
	"gopkg.in/validator.v2"
)

type KSTestConf struct {
	Strictness      uint   `yaml:"strictness" validate:"min=1,max=5" default:"1"`
	NormalCount     uint   `yaml:"normal_count" validate:"nonzero" default:"3"`
	MaxFingerprints uint   `yaml:"max_fingerprints" validate:"nonzero" default:"20"`
}

func (self *KSTestConf) Default() {
	err := self.Validate()
	if err != nil {
		errs := err.(validator.ErrorMap)
		for fieldName := range errs {
			log.Infof("Using default for KSTestConf.%s", fieldName)
			defaults.SetDefault(self, fieldName)
		}
	}

}

func (self *KSTestConf) Validate() error {
	return validator.Validate(self)
}

