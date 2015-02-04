package grafana

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	"github.com/nvcook42/morgoth/defaults"
	"github.com/nvcook42/morgoth/engine"
)

type GrafanaConf struct {
	URL        string            `yaml:"url" validate:"nonzero" default:"http://grafanarel.s3.amazonaws.com/grafana-1.9.1.tar.gz"`
	Port       uint              `yaml:"port" validate:"min=1,max=65535" default:"8080"`
	Dir        string            `yaml:"dir" validate:"nonzero" default:"grafana_tmp"`
	EngineConf engine.EngineConf `yaml:"engine"`
	GrafanaDB  string            `yaml:"grafana_db" validate:"nonzero" default:"grafana"`
}

func (self *GrafanaConf) Validate() error {
	valid := self.EngineConf.Validate()
	if valid != nil {
		return valid
	}
	return validator.Validate(self)
}

//Sets any invalid fields to their default value
func (self *GrafanaConf) Default() {
	self.EngineConf.Default()
	err := self.Validate()
	if err != nil {
		errs := err.(validator.ErrorMap)
		for fieldName := range errs {
			if ok, _ := defaults.HasDefault(self, fieldName); ok {
				glog.Infof("Using default for GrafanaConf.%s", fieldName)
				defaults.SetDefault(self, fieldName)
			}
		}
	}
}
