package config

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/engine"
	"github.com/nvcook42/morgoth/metric"
	"github.com/nvcook42/morgoth/fitting"
)

// Base config struct for the entire morgoth config
type Config struct {
	EngineConf engine.EngineConf      `yaml:"data_engine"`
	Metrics    []metric.MetricConf    `yaml:"metrics"`
	Fittings   []fitting.FittingConf  `yaml:"fittings"`
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
