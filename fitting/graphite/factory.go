package graphite

import (
	"errors"
	"fmt"
	"github.com/nathanielc/morgoth/config/types"
	"github.com/nathanielc/morgoth/fitting"
)

type GraphiteFactory struct {
}

func (self *GraphiteFactory) NewConf() types.Configuration {
	return new(GraphiteConf)
}

func (self *GraphiteFactory) GetInstance(config types.Configuration) (interface{}, error) {
	graphiteConf, ok := config.(*GraphiteConf)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Configuration is not GraphiteConf %v", config))
	}
	return &GraphiteFitting{port: graphiteConf.Port}, nil
}

func init() {
	factory := new(GraphiteFactory)
	fitting.Registery.RegisterFactory("graphite", factory)
}
