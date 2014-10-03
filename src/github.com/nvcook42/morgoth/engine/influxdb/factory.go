package influxdb

import (
	"github.com/nvcook42/morgoth/engine"
	"github.com/nvcook42/morgoth/registery"
)

type InfluxDBFactory struct {
}

func (self *InfluxDBFactory) NewConf() registery.Configuration {
	return new(InfluxDBConf)
}

func (self *InfluxDBFactory) GetInstance(config registery.Configuration) (interface{}, error) {
	return new(InfluxDBEngine), nil
}

func init() {
	factory := new(InfluxDBFactory)
	engine.Registery.RegisterFactory("influxdb", factory)
}
