package influxdb

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	config "github.com/nvcook42/morgoth/config/types"
)

type InfluxDBConf struct {
	Host     string `validate:"min=1"           default:"localhost"`
	Port     uint   `validate:"min=1,max=65535" default:"8083"`
	User     string `validate:"min=1"`
	Password string `validate:"min=1"`
	Database string `validate:"min=1"`
}

func (self *InfluxDBConf) Validate() error {
	return validator.Validate(self)
}

func (self *InfluxDBConf) Default() {
	config.PerformDefault(self)
}
