package config

import (
	"errors"
	"github.com/nvcook42/morgoth/detectors/mgof"
	"github.com/nvcook42/morgoth/validator"
)

type DetectorType string

const (
	MGOF DetectorType = "mgof"
)

func (self DetectorType) Validate() error {
	switch self {
	case
		MGOF:
		return nil
	}
	return errors.New("Invalid DetectorType")
}

type Detector struct {
	Type DetectorType   `yaml:"type"`
	MGOF *mgof.MGOFConf `yaml:"mgof,omitempty"`
}

func (self *Detector) Default() {
	switch self.Type {
	case MGOF:
		self.MGOF.Default()
	}
}

func (self Detector) Validate() error {
	return validator.ValidateAll(self)
}
