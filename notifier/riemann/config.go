package riemann

import (
	"github.com/golang/glog"
	"github.com/nvcook42/morgoth/defaults"
	"gopkg.in/validator.v2"
)

type RiemannConf struct {
	Host string `validate:"min=1" default:"localhost"`
	Port uint `validate:"min=1,max=65535" default:"5555"`
}

func (self *RiemannConf) Validate() error {
	return validator.Validate(self)
}

//Sets any invalid fields to their default value
func (self *RiemannConf) Default() {
	err := self.Validate()
	if err != nil {
		errs := err.(validator.ErrorMap)
		for fieldName := range errs {
			if ok, _ := defaults.HasDefault(self, fieldName); ok {
				glog.Infof("Using default for RiemannConf.%s", fieldName)
				defaults.SetDefault(self, fieldName)
			}
		}
	}
}
