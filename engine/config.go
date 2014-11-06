package engine

import (
	//"github.com/golang/glog"
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

func FromYAML(yaml string) (*EngineConf, error) {
	conf := new(EngineConf)
	err := dynamic_type.PerformFromYAML(yaml, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
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
