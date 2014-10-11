package engine

import (
	//log "github.com/cihub/seelog"
	"errors"
	"fmt"
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
	instance, err := self.PerformGetInstance(Registery)
	if err != nil {
		return nil, err
	}

	engine, ok := instance.(Engine)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Instance %v is not of type Engine", instance))
	}
	return engine, nil
}
