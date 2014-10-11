package notifier

import (
	"github.com/nvcook42/morgoth/config/dynamic_type"
)

type NotifierConf struct {
	dynamic_type.DynamicConfiguration
}

func (self *NotifierConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return self.PerformUnmarshalYAML(Registery, unmarshal)
}
