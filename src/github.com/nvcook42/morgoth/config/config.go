package config

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/engine"
)

// Base config struct for the entire morgoth config
type Config struct {
	DataEngine engine.DataEngine `yaml:"data_engine"`
	Metrics    []Metric          `yaml:"metrics"`
	Fittings   []Fitting         `yaml:"fittings"`
}

func (self *Config) Default() {
	self.DataEngine.Default()
	for i := range self.Metrics {
		self.Metrics[i].Default()
	}
	for i := range self.Fittings {
		self.Fittings[i].Default()
	}
}

func (self Config) Validate() error {
	log.Debugf("Validating Config %v", self)
	valid := self.DataEngine.Validate()
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
