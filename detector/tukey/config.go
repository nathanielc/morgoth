package tukey

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	"github.com/nvcook42/morgoth/defaults"
)

type TukeyConf struct {
	Threshold float64 `yaml:"threshold" validate:"min=0,nonzero" default:"3"`
}

func (self *TukeyConf) Default() {
	err := self.Validate()
	if err != nil {
		errs, ok := err.(validator.ErrorMap)
		if !ok {
			// Non validation error returned can't continue
			return
		}
		for fieldName := range errs {
			glog.Infof("Using default for TukeyConf.%s", fieldName)
			defaults.SetDefault(self, fieldName)
		}
	}

}

func (self *TukeyConf) Validate() error {
	return validator.Validate(self)
}
