package influxdb

import (
	"github.com/nvcook42/morgoth/engine"
)

type InfluxDBEngine struct {
}

func (self *InfluxDBEngine) Initialize() error {
	return nil
}

func (self *InfluxDBEngine) GetReader() engine.Reader {
	return nil
}

func (self *InfluxDBEngine) GetWriter() engine.Writer {
	return nil
}
