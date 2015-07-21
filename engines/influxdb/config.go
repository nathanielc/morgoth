package influxdb

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	"github.com/nathanielc/morgoth/config"
)

type InfluxDBConf struct {
	Host               string `validate:"nonzero"           default:"localhost"`
	Port               uint   `validate:"min=1,max=65535" default:"8083"`
	User               string `validate:"nonzero"`
	Password           string `validate:"nonzero"`
	Database           string `yaml:"database" validate:"nonzero"`
	AnomalyMeasurement string `yaml:"anomaly_measurement" validate:"nonzero"`
	MeasurementTag     string `yaml:"measurement_tag" validate:"nonzero"`
}

func (self *InfluxDBConf) Validate() error {
	return validator.Validate(self)
}

func (self *InfluxDBConf) Default() {
	config.PerformDefault(self)
}
