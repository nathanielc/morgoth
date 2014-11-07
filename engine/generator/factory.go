package generator

import (
	"errors"
	"fmt"
	"github.com/nvcook42/morgoth/config/types"
	"github.com/nvcook42/morgoth/engine"
)

type GeneratorFactory struct {
}

func (self *GeneratorFactory) NewConf() types.Configuration {
	return new(GeneratorConf)
}

func (self *GeneratorFactory) GetInstance(config types.Configuration) (interface{}, error) {
	conf, ok := config.(*GeneratorConf)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Configuration is not GeneratorConf%v", config))
	}
	engine := &GeneratorEngine{
		config: conf,
	}
	return engine, nil
}

func init() {
	factory := new(GeneratorFactory)
	engine.Registery.RegisterFactory("generator", factory)
}
