package detector

import (
	"github.com/nvcook42/morgoth/config/dynamic_type"
)

type DetectorConf struct {
	dynamic_type.DynamicConfiguration
}

func (self *DetectorConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return self.PerformUnmarshalYAML(Registery, unmarshal)
}
