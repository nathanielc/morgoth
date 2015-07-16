package influxdb

import (
	"errors"
	"fmt"
	"github.com/nathanielc/morgoth/config/types"
	"github.com/nathanielc/morgoth/engine"
)

type InfluxDBFactory struct {
}

func (self *InfluxDBFactory) NewConf() types.Configuration {
	return new(InfluxDBConf)
}

func (self *InfluxDBFactory) GetInstance(config types.Configuration) (interface{}, error) {
	conf, ok := config.(*InfluxDBConf)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Configuration is not InfluxDBConf%v", config))
	}
	engine := &InfluxDBEngine{
		config: conf,
	}
	return engine, nil
}

func init() {
	factory := new(InfluxDBFactory)
	engine.Registery.RegisterFactory("influxdb", factory)
}
