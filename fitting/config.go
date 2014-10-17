package fitting

import (
	"errors"
	"fmt"
	"github.com/nvcook42/morgoth/config/dynamic_type"
)

type FittingConf struct {
	dynamic_type.DynamicConfiguration
}

func (self *FittingConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return self.PerformUnmarshalYAML(Registery, unmarshal)
}

func (self *FittingConf) GetFitting() (Fitting, error) {
	instance, err := self.PerformGetInstance(Registery)
	if err != nil {
		return nil, err
	}

	fitting, ok := instance.(Fitting)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Instance %v is not of type Fitting", instance))
	}
	return fitting, nil

}
