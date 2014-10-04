package influxdb

import (
	"github.com/nvcook42/morgoth/engine"
	"github.com/nvcook42/morgoth/config/types"
)

type InfluxDBFactory struct {
}

func (self *InfluxDBFactory) NewConf() types.Configuration {
	return new(InfluxDBConf)
}

func (self *InfluxDBFactory) GetInstance(config types.Configuration) (interface{}, error) {
	return new(InfluxDBEngine), nil
}

func init() {
	factory := new(InfluxDBFactory)
	engine.Registery.RegisterFactory("influxdb", factory)
}
