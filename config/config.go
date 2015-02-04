package config

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	app "github.com/nvcook42/morgoth/app/types"
	"github.com/nvcook42/morgoth/engine"
	"github.com/nvcook42/morgoth/fitting"
	"github.com/nvcook42/morgoth/metric"
	"github.com/nvcook42/morgoth/schedule"
)

// Base config struct for the entire morgoth config
type Config struct {
	EngineConf engine.EngineConf             `yaml:"engine"`
	Metrics    []metric.MetricSupervisorConf `yaml:"metrics"`
	Fittings   []fitting.FittingConf         `yaml:"fittings"`
	Schedule   schedule.ScheduleConf         `yaml:"schedule"`
	Morgoth    MorgothConf                   `yaml:"morgoth"`
}

func (self *Config) Default() {
	self.EngineConf.Default()
	self.Schedule.Default()
	self.Morgoth.Default()
	for i := range self.Metrics {
		self.Metrics[i].Default()
	}
	for i := range self.Fittings {
		self.Fittings[i].Default()
	}
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
		if err == nil && fitting != nil {
			fittings[i] = fitting
		} else {
			glog.Error(err)
		}
	}

	return fittings
}

func (self *Config) GetSchedule() schedule.Schedule {
	return self.Schedule.GetSchedule()
}
