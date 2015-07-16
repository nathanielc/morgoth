package config

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	app "github.com/nathanielc/morgoth/app/types"
	"github.com/nathanielc/morgoth/engine"
	"github.com/nathanielc/morgoth/fitting"
	"github.com/nathanielc/morgoth/metric"
	"github.com/nathanielc/morgoth/schedule"
)

// Base config struct for the entire morgoth config
type Config struct {
	EngineConf engine.EngineConf     `yaml:"engine"`
	Schedule   schedule.ScheduleConf `yaml:"schedule"`
	Morgoth    MorgothConf           `yaml:"morgoth"`
}

func (self *Config) Default() {
	self.EngineConf.Default()
	self.Schedule.Default()
	self.Morgoth.Default()
}

func (self Config) Validate() error {
	glog.V(2).Infof("Validating Config %v", self)
	valid := self.EngineConf.Validate()
	if valid != nil {
		return valid
	}

	valid = self.Schedule.Validate()
	if valid != nil {
		return valid
	}

	valid = self.Morgoth.Validate()
	if valid != nil {
		return valid
	}
	return nil
}

func (self *Config) GetSchedule() schedule.Schedule {
	return self.Schedule.GetSchedule()
}
