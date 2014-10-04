package engine

import (
	//log "github.com/cihub/seelog"
	"errors"
	"github.com/nvcook42/morgoth/registery"
	"github.com/nvcook42/morgoth/config/dynamic_type"
)


// The Data Engine subsection of the config
type EngineConf struct {
	Type string
	Conf registery.Configuration
}

func (self *EngineConf) Default() {
	if self.Conf != nil {
		self.Conf.Default()
	}
}

func (self EngineConf) Validate() error {
	if self.Conf == nil {
		return errors.New("No conf found")
	}
	return self.Conf.Validate()
}

func (self *EngineConf) GetEngine() (Engine, error) {
	factory, err := Registery.GetFactory(self.Type)
	if err != nil {
		return nil, err
	}
	engine, err := factory.GetInstance(self.Conf)
	return engine.(Engine), err
}

func (self *EngineConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	engineType, config, err := dynamic_type.UnmarshalDynamicType(Registery, unmarshal)
	if err != nil {
		return err
	}
	self.Type = engineType
	self.Conf = config
	return nil
}

