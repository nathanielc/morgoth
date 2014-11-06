package influxdb

import (
	"github.com/golang/glog"
	"github.com/nvcook42/morgoth/defaults"
	"gopkg.in/validator.v2"
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

//Sets any invalid fields to their default value
func (self *InfluxDBConf) Default() {
	err := self.Validate()
	if err != nil {
		errs := err.(validator.ErrorMap)
		for fieldName := range errs {
			if ok, _ := defaults.HasDefault(self, fieldName); ok {
				glog.Infof("Using default for InfluxDBConf.%s", fieldName)
				defaults.SetDefault(self, fieldName)
			}
		}
	}
}
