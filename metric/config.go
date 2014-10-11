package metric

import (
	log "github.com/cihub/seelog"
	"errors"
	app "github.com/nvcook42/morgoth/app/types"
	"github.com/nvcook42/morgoth/detector"
	"github.com/nvcook42/morgoth/notifier"
	"github.com/nvcook42/morgoth/schedule"
	"regexp"
)

// Represents a single metric conf section
type MetricConf struct {
	Pattern   Pattern                 `yaml:"pattern"`
	Schedule  schedule.ScheduleConf   `yaml:"schedule"`
	Detectors []detector.DetectorConf `yaml:"detectors"`
	Notifiers []notifier.NotifierConf `yaml:"notifiers"`
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

func (self *MetricConf) GetSupervisor(app app.App) Supervisor {

	detectors := make([]detector.Detector, 0, len(self.Detectors))
	for i := range self.Detectors {
		detector, err := self.Detectors[i].GetDetector()
		if err == nil {
			detectors = append(detectors, detector)
		} else {
			log.Errorf("Error getting configured detector: %s", err.Error())
		}
	}

	notifiers := make([]notifier.Notifier, 0, len(self.Notifiers))
	for i := range self.Notifiers {
		notifier, err := self.Notifiers[i].GetNotifier()
		if err == nil {
			notifiers = append(notifiers, notifier)
		} else {
			log.Errorf("Error getting configured notifier: %s", err.Error())
		}
	}

	return NewSupervisor(
		self.Pattern,
		app.GetWriter(),
		detectors,
		notifiers,
		self.Schedule.GetSchedule(),
	)
}

type Pattern string

func (self Pattern) Validate() error {
	if len(self) == 0 {
		return errors.New("Pattern cannot be empty")
	}
	_, err := regexp.Compile(string(self))
	return err
}
