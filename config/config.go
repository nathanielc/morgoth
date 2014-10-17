package config

import (
	log "github.com/cihub/seelog"
	app "github.com/nvcook42/morgoth/app/types"
	"github.com/nvcook42/morgoth/engine"
	"github.com/nvcook42/morgoth/fitting"
	"github.com/nvcook42/morgoth/metric"
)

// Base config struct for the entire morgoth config
type Config struct {
	EngineConf engine.EngineConf     `yaml:"engine"`
	Metrics    []metric.MetricConf   `yaml:"metrics"`
	Fittings   []fitting.FittingConf `yaml:"fittings"`
}

func (self *Config) Default() {
	self.EngineConf.Default()
	for i := range self.Metrics {
		self.Metrics[i].Default()
	}
	for i := range self.Fittings {
		self.Fittings[i].Default()
	}
}

func (self Config) Validate() error {
	log.Debugf("Validating Config %v", self)
	valid := self.EngineConf.Validate()
	if valid != nil {
		return valid
	}

	for i := range self.Metrics {
		valid := self.Metrics[i].Validate()
		if valid != nil {
			return valid
		}
	}

	for i := range self.Fittings {
		valid := self.Fittings[i].Validate()
		if valid != nil {
			return valid
		}
	}
	return nil
}

func (self *Config) GetSupervisors(app app.App) []metric.Supervisor {
	supervisors := make([]metric.Supervisor, len(self.Metrics))
	for i := range self.Metrics {
		supervisors[i] = self.Metrics[i].GetSupervisor(app)
	}

	return supervisors
}

func (self *Config) GetFittings() []fitting.Fitting {
	fittings := make([]fitting.Fitting, len(self.Fittings))
	for i := range self.Fittings {
		fitting, err := self.Fittings[i].GetFitting()
		if err == nil {
			fittings[i] = fitting
		}
	}

	return fittings
}
