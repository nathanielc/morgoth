package engine

import (
	//log "github.com/cihub/seelog"
	"errors"
	"github.com/nvcook42/morgoth/registery"
	"github.com/nvcook42/morgoth/config/dynamic_type"
)


// The Data Engine subsection of the config
type DataEngine struct {
	Type string
	Conf registery.Configuration
}

func (self *DataEngine) Default() {
	if self.Conf != nil {
		self.Conf.Default()
	}
}

func (self DataEngine) Validate() error {
	if self.Conf == nil {
		return errors.New("No conf found")
	}
	return self.Conf.Validate()
}

func (self *DataEngine) GetEngine() (Engine, error) {
	factory, err := Registery.GetFactory(self.Type)
	if err != nil {
		return nil, err
	}
	engine, err := factory.GetInstance(self.Conf)
	return engine.(Engine), err
}

func (self *DataEngine) UnmarshalYAML(unmarshal func(interface{}) error) error {
	engineType, config, err := dynamic_type.UnmarshalDynamicType("engine", Registery, unmarshal)
	if err != nil {
		return err
	}
	self.Type = engineType
	self.Conf = config
	return nil
}

