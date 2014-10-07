package fitting

import (
	"github.com/nvcook42/morgoth/config/dynamic_type"
)

type FittingConf struct {
	dynamic_type.DynamicConfiguration
}

func (self *FittingConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return self.PerformUnmarshalYAML(Registery, unmarshal)
}
