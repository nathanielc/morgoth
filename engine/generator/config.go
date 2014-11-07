package generator

import (
	"github.com/golang/glog"
	"github.com/nvcook42/morgoth/defaults"
	"gopkg.in/validator.v2"
)

type GeneratorConf struct {
}

func (self *GeneratorConf) Validate() error {
	return validator.Validate(self)
}

//Sets any invalid fields to their default value
func (self *GeneratorConf) Default() {
	err := self.Validate()
	if err != nil {
		errs := err.(validator.ErrorMap)
		for fieldName := range errs {
			if ok, _ := defaults.HasDefault(self, fieldName); ok {
				glog.Infof("Using default for GeneratorConf.%s", fieldName)
				defaults.SetDefault(self, fieldName)
			}
		}
	}
}
