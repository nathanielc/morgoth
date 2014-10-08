package metric

import (
	"errors"
	"github.com/nvcook42/morgoth/detector"
	"github.com/nvcook42/morgoth/schedule"
	"regexp"
)

// Represents a single metric conf section
type MetricConf struct {
	Pattern   Pattern                 `yaml:"pattern"`
	Schedule  schedule.ScheduleConf   `yaml:"schedule"`
	Detectors []detector.DetectorConf `yaml:"detectors"`
}

func (self *MetricConf) Default() {
	self.Schedule.Default()

	for i := range self.Detectors {
		self.Detectors[i].Default()
	}
}

func (self MetricConf) Validate() error {
	if valid := self.Pattern.Validate(); valid != nil {
		return valid
	}
	if valid := self.Schedule.Validate(); valid != nil {
		return valid
	}
	for i := range self.Detectors {
		if valid := self.Detectors[i].Validate(); valid != nil {
			return valid
		}
	}
	return nil
}

type Pattern string

func (self Pattern) Validate() error {
	if len(self) == 0 {
		return errors.New("Pattern cannot be empty")
	}
	_, err := regexp.Compile(string(self))
	return err
}
