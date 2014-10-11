package notifier

import (
	"errors"
	"fmt"
	"github.com/nvcook42/morgoth/config/dynamic_type"
)

type NotifierConf struct {
	dynamic_type.DynamicConfiguration
}

func (self *NotifierConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return self.PerformUnmarshalYAML(Registery, unmarshal)
}

func (self *NotifierConf) GetNotifier() (Notifier, error) {
	instance, err := self.PerformGetInstance(Registery)
	if err != nil {
		return nil, err
	}

	notifier, ok := instance.(Notifier)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Instance %v is not of type Notifier", instance))
	}
	return notifier, nil
}
