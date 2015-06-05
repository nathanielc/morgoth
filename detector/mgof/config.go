package mgof

import (
	"errors"
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	config "github.com/nvcook42/morgoth/config/types"
)

type MGOFConf struct {
	Min float64 `yaml:"min"`
	Max float64 `yaml:"max"`
	// A smaller number means that the algorithm is more distinguishing
	NullConfidence  uint `yaml:"null_confidence" validate:"min=1,max=10" default:"10"`
	NBins           uint `yaml:"nbins" validate:"nonzero" default:"15"`
	NormalCount     uint `yaml:"normal_count" validate:"nonzero" default:"2"`
	MaxFingerprints uint `yaml:"max_fingerprints" validate:"nonzero" default:"20"`
}

func (self *MGOFConf) Default() {
	config.PerformDefault(self)
}

func (self *MGOFConf) Validate() error {
	if self.Min >= self.Max {
		return errors.New("Min must be less that Max")
	}
	return validator.Validate(self)
}
