package metric

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/detector"
	"github.com/nvcook42/morgoth/engine"
	"github.com/nvcook42/morgoth/metric/set"
	"github.com/nvcook42/morgoth/metric/types"
	"github.com/nvcook42/morgoth/notifier"
	"github.com/nvcook42/morgoth/schedule"
	"time"
)

// A supervisor keeps track of all metrics that match its pattern
// proxies Detect calls to the associated detectors
// It also supervises the notification of anomalous metrics
type Supervisor interface {
	GetPattern() types.Pattern
	AddMetric(types.MetricID)
	Detect(rotation *schedule.Rotation, start time.Time, stop time.Time)
}

type SupervisorStruct struct {
	pattern   types.Pattern
	writer    engine.Writer
	detectors map[schedule.Rotation][]detector.Detector
	notifiers []notifier.Notifier
	metrics   *set.Set
}

func NewSupervisor(
	pattern types.Pattern,
	writer engine.Writer,
	detectors map[schedule.Rotation][]detector.Detector,
	notifiers []notifier.Notifier,
) *SupervisorStruct {

	s := &SupervisorStruct{
		pattern:   pattern,
		writer:    writer,
		detectors: detectors,
		notifiers: notifiers,
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

func (self *SupervisorStruct) Detect(rotation *schedule.Rotation, start time.Time, stop time.Time) {
	detectors := self.detectors[*rotation]
	self.metrics.Each(func(metric types.MetricID) {
		for _, detector := range detectors {
			if detector.Detect(metric, start, stop) {
				log.Infof("Metric %s is anomalous", metric)
			}
		}
	})
}
