package grafana

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	config "github.com/nathanielc/morgoth/config/types"
	"github.com/nathanielc/morgoth/engine/influxdb"
)

type GrafanaConf struct {
	URL          string                `yaml:"url"  validate:"nonzero" default:"http://grafanarel.s3.amazonaws.com/grafana-1.9.1.tar.gz"`
	Port         uint                  `yaml:"port" validate:"min=1,max=65535" default:"8080"`
	Dir          string                `yaml:"dir"  validate:"nonzero" default:"grafana_tmp"`
	InfluxDBConf influxdb.InfluxDBConf `yaml:"influxdb"`
	GrafanaDB    string                `yaml:"grafana_db" validate:"nonzero" default:"grafana"`
}

func (self *GrafanaConf) Validate() error {
	valid := self.InfluxDBConf.Validate()
	if valid != nil {
		return valid
	}
	return validator.Validate(self)
}

//Sets any invalid fields to their default value
func (self *GrafanaConf) Default() {
	self.InfluxDBConf.Default()
	config.PerformDefault(self)
}
