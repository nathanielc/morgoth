package morgoth

import (
	"errors"
	"fmt"
	"github.com/nathanielc/morgoth/config"
)

type Engine interface {
	Initialize() error
	GetWindows(query Query) ([]*Window, error)
	NewQueryBuilder(queryTemplate string) (QueryBuilder, error)
	RecordAnomalous(w Window) error
}

var EngineRegistery *config.Registery

func init() {
	EngineRegistery = config.NewRegistry()
}

// Configuration

// The Data Engine subsection of the config
type EngineConf struct {
	config.DynamicConfiguration
}

func (self *EngineConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return self.PerformUnmarshalYAML(EngineRegistery, unmarshal)
}

func EngineFromYAML(yaml string) (*EngineConf, error) {
	conf := new(EngineConf)
	err := config.PerformFromYAML(yaml, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (self *EngineConf) GetEngine() (Engine, error) {
	instance, err := self.PerformGetInstance(EngineRegistery)
	if err != nil {
		return nil, err
	}

	engine, ok := instance.(Engine)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Instance %v is not of type Engine", instance))
	}
	return engine, nil
}
