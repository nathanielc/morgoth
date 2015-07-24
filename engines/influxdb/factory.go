package influxdb

import (
	"errors"
	"fmt"
	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/influxdb/influxdb/client"
	"github.com/nathanielc/morgoth/config"
	"net/url"
)

type InfluxDBFactory struct {
}

func (self *InfluxDBFactory) NewConf() config.Configuration {
	return new(InfluxDBConf)
}

func (self *InfluxDBFactory) GetInstance(config config.Configuration) (interface{}, error) {
	conf, ok := config.(*InfluxDBConf)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Configuration is not InfluxDBConf%v", config))
	}
	url, err := url.Parse(fmt.Sprintf("http://%s:%d", conf.Host, conf.Port))
	if err != nil {
		return nil, err
	}
	engine := &InfluxDBEngine{
		conf: client.Config{
			URL:      *url,
			Username: conf.User,
			Password: conf.Password,
		},
		database:           conf.Database,
		anomalyMeasurement: conf.AnomalyMeasurement,
		measurementTag:     conf.MeasurementTag,
	}
	return engine, nil
}

func init() {
	factory := new(InfluxDBFactory)
	morgoth.EngineRegistery.RegisterFactory("influxdb", factory)
}
