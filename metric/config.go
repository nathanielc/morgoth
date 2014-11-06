package metric

import (
	"github.com/golang/glog"
	app "github.com/nvcook42/morgoth/app/types"
	"github.com/nvcook42/morgoth/detector"
	"github.com/nvcook42/morgoth/metric/types"
	"github.com/nvcook42/morgoth/notifier"
	"github.com/nvcook42/morgoth/schedule"
)

// Represents a single metric conf section
type MetricSupervisorConf struct {
	Pattern   types.Pattern           `yaml:"pattern"`
	Detectors []detector.DetectorConf `yaml:"detectors"`
	Notifiers []notifier.NotifierConf `yaml:"notifiers"`
}

func (self *MetricSupervisorConf) Default() {
	for i := range self.Detectors {
		self.Detectors[i].Default()
	}
}

func (self MetricSupervisorConf) Validate() error {
	if valid := self.Pattern.Validate(); valid != nil {
		return valid
	}
	for i := range self.Detectors {
		if valid := self.Detectors[i].Validate(); valid != nil {
			return valid
		}
	}
	return nil
}

func (self *MetricSupervisorConf) GetSupervisor(app app.App) Supervisor {

	schd := app.GetSchedule()
	detectorsMap := make(map[schedule.Rotation][]detector.Detector, len(schd.Rotations))
	for _, rotation := range schd.Rotations {
		detectors := make([]detector.Detector, 0, len(self.Detectors))
		for i := range self.Detectors {
			detector, err := self.Detectors[i].GetDetector()
			if err == nil {
				err = detector.Initialize(app, rotation)
				if err == nil {
					detectors = append(detectors, detector)
				} else {
					glog.Error("Error initializing detector ", err)
				}
			} else {
				glog.Errorf("Error getting configured detector: %s", err.Error())
			}
		}
		detectorsMap[rotation] = detectors
	}

	notifiers := make([]notifier.Notifier, 0, len(self.Notifiers))
	for i := range self.Notifiers {
		notifier, err := self.Notifiers[i].GetNotifier()
		if err == nil {
			notifiers = append(notifiers, notifier)
		} else {
			glog.Errorf("Error getting configured notifier: %s", err.Error())
		}
	}

	return NewSupervisor(
		self.Pattern,
		app.GetWriter(),
		detectorsMap,
		notifiers,
	)
}
