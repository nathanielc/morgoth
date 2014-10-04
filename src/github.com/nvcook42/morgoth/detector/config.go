
package detector

import (
	"errors"
	"github.com/nvcook42/morgoth/config/types"
	"github.com/nvcook42/morgoth/config/dynamic_type"
)


type DetectorConf struct {
	Type string
	Conf types.Configuration
}

func (self *DetectorConf) Default() {
	if self.Conf != nil {
		self.Conf.Default()
	}
}

func (self DetectorConf) Validate() error {
	if self.Conf == nil {
		return errors.New("No detector conf found")
	}
	return self.Conf.Validate()
}

func (self *DetectorConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	detectorType, config, err := dynamic_type.UnmarshalDynamicType(Registery, unmarshal)
	if err != nil {
		return err
	}
	self.Type = detectorType
	self.Conf = config
	return nil
}
