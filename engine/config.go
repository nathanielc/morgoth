package engine

import (
	//log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/config/dynamic_type"
)

// The Data Engine subsection of the config
type EngineConf struct {
	dynamic_type.DynamicConfiguration
}

func (self *EngineConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return self.PerformUnmarshalYAML(Registery, unmarshal)
}

func (self *EngineConf) GetEngine() (Engine, error) {
	factory, err := Registery.GetFactory(self.Type)
	if err != nil {
		return nil, err
	}
	engine, err := factory.GetInstance(self.Conf)
	return engine.(Engine), err
}
