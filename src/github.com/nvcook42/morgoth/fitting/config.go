package fitting

import (
	"errors"
	"github.com/nvcook42/morgoth/registery"
	"github.com/nvcook42/morgoth/config/dynamic_type"
)

type FittingConf struct {
	Type string
	Conf registery.Configuration
}

func (self *FittingConf) Default() {
	if self.Conf != nil {
		self.Conf.Default()
	}
}

func (self FittingConf) Validate() error {
	if self.Conf == nil {
		return errors.New("No fitting conf found")
	}
	return self.Conf.Validate()
}


func (self *FittingConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	fittingType, config, err := dynamic_type.UnmarshalDynamicType(Registery, unmarshal)
	if err != nil {
		return err
	}
	self.Type = fittingType
	self.Conf = config
	return nil
}
