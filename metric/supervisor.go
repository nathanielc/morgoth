package metric

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/detector"
	"github.com/nvcook42/morgoth/engine"
	"github.com/nvcook42/morgoth/metric/set"
	"github.com/nvcook42/morgoth/metric/types"
	app "github.com/nvcook42/morgoth/app/types"
	"github.com/nvcook42/morgoth/notifier"
	"github.com/nvcook42/morgoth/schedule"
	"time"
)

// A supervisors schedules all the actions
// associated with detecting, notifing etc metrics
// that match their pattern
type Supervisor interface {
	GetPattern() types.Pattern
	AddMetric(types.MetricID)
	Start(app.App)
}

type SupervisorStruct struct {
	pattern   types.Pattern
	writer    engine.Writer
	detectors []detector.Detector
	notifiers []notifier.Notifier
	schedule  schedule.Schedule
	metrics   *set.Set
}

func NewSupervisor(
	pattern types.Pattern,
	writer engine.Writer,
	detectors []detector.Detector,
	notifiers []notifier.Notifier,
	schedule schedule.Schedule,
) *SupervisorStruct {

	s := &SupervisorStruct{
		pattern:   pattern,
		writer:    writer,
		detectors: detectors,
		notifiers: notifiers,
		schedule:  schedule,
		metrics:   set.New(0),
	}
	return s
}

func (self *SupervisorStruct) GetPattern() types.Pattern {
	return self.pattern
}

func (self *SupervisorStruct) AddMetric(metric types.MetricID) {
	self.metrics.Add(metric)
}

func (self *SupervisorStruct) Start(app app.App) {
	self.schedule.Callback = self.detect
	self.schedule.Start()

	for _, detector := range self.detectors {
		err := detector.Initialize(app)
		if err != nil {
			log.Warnf("Failed to initial detector %v %s", detector, err.Error())
		}
	}
}

func (self *SupervisorStruct) GetSchedule() schedule.Schedule {
	return self.schedule
}

func (self *SupervisorStruct) detect(start time.Time, stop time.Time) {
	self.metrics.Each(func(metric types.MetricID) {
		for _, detector := range self.detectors {
			if detector.Detect(metric, start, stop) {
				log.Infof("Metric %s is anomalous", metric)
			}
		}
	})
}
