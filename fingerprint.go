package morgoth

import (
	"errors"
	"fmt"
	"github.com/nathanielc/morgoth/config"
	"github.com/nathanielc/morgoth/counter"
)

type Fingerprint interface {
	IsMatch(other counter.Countable) bool
}

type Fingerprinter interface {
	Fingerprint(window *Window) Fingerprint
}

var FingerprinterRegistery *config.Registery

func init() {
	FingerprinterRegistery = config.NewRegistry()
}

type FingerprinterConf struct {
	config.DynamicConfiguration
}

func (self *FingerprinterConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return self.PerformUnmarshalYAML(FingerprinterRegistery, unmarshal)
}

func FingerprinterFromYAML(yaml string) (*FingerprinterConf, error) {
	conf := new(FingerprinterConf)
	err := config.PerformFromYAML(yaml, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (self *FingerprinterConf) GetFingerprinter() (Fingerprinter, error) {
	instance, err := self.PerformGetInstance(FingerprinterRegistery)
	if err != nil {
		return nil, err
	}

	fingerprinter, ok := instance.(Fingerprinter)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Instance %v is not of type Fingerprinter", instance))
	}
	return fingerprinter, nil
}
