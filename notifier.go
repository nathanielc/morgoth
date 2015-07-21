package morgoth

import (
	"errors"
	"fmt"
	"github.com/nathanielc/morgoth/config"
)

type Notifier interface {
	Notify(msg string, w *Window)
}

var (
	NotifierRegistery *config.Registery
)

func init() {
	NotifierRegistery = config.NewRegistry()
}

type NotifierConf struct {
	config.DynamicConfiguration
}

func (self *NotifierConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return self.PerformUnmarshalYAML(NotifierRegistery, unmarshal)
}

func (self *NotifierConf) GetNotifier() (Notifier, error) {
	instance, err := self.PerformGetInstance(NotifierRegistery)
	if err != nil {
		return nil, err
	}

	notifier, ok := instance.(Notifier)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Instance %v is not of type Notifier", instance))
	}
	return notifier, nil
}
