package mgof

import (
	"errors"
	"github.com/golang/glog"
	"github.com/nvcook42/morgoth/defaults"
	"gopkg.in/validator.v2"
)

type MGOFConf struct {
	Min             float64 `yaml:"min"`
	Max             float64 `yaml:"max"`
	CHI2            float64 `yaml:"chi2" validate:"min=0,max=1,nonzero" default:"0.95"`
	NBins           uint    `yaml:"nbins" validate:"nonzero" default:"15"`
	NormalCount     uint    `yaml:"normal_count" validate:"nonzero" default:"3"`
	MaxFingerprints uint    `yaml:"max_fingerprints" validate:"nonzero" default:"20"`
}

func (self *MGOFConf) Default() {
	err := self.Validate()
	if err != nil {
		errs, ok := err.(validator.ErrorMap)
		if !ok {
			// Non validation error returned can't continue
			return
		}
		for fieldName := range errs {
			glog.Infof("Using default for MGOFConf.%s", fieldName)
			defaults.SetDefault(self, fieldName)
		}
	}

}

func (self *MGOFConf) Validate() error {
	if self.Min > self.Max {
		return errors.New("Min must be less that Max")
	}
	return validator.Validate(self)
}
