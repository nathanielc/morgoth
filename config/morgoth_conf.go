package config

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	"github.com/nvcook42/morgoth/defaults"
)

type MorgothConf struct {
	MetaDir string      `yaml:"meta_dir" validate:"nonzero" default:"meta"`
}

func (self MorgothConf) Validate() error {
	return validator.Validate(self)
}

func (self *MorgothConf) Default() {
	err := self.Validate()
	if err != nil {
		errs := err.(validator.ErrorMap)
		for fieldName := range errs {
			if ok, _ := defaults.HasDefault(self, fieldName); ok {
				glog.Infof("Using default for Morgoth.%s", fieldName)
				defaults.SetDefault(self, fieldName)
			}
		}
	}
}
