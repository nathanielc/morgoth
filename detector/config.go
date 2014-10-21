package detector

import (
	"errors"
	"fmt"
	"github.com/nvcook42/morgoth/config/dynamic_type"
)

type DetectorConf struct {
	dynamic_type.DynamicConfiguration
}

func (self *DetectorConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return self.PerformUnmarshalYAML(Registery, unmarshal)
}

func FromYAML(yaml string) (*DetectorConf, error) {
	conf := new(DetectorConf)
	err := dynamic_type.PerformFromYAML(yaml, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (self *DetectorConf) GetDetector() (Detector, error) {
	instance, err := self.PerformGetInstance(Registery)
	if err != nil {
		return nil, err
	}

	detector, ok := instance.(Detector)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Instance %v is not of type Detector", instance))
	}
	return detector, nil
}
