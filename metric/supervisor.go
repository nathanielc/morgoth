package metric

import (
	"github.com/nvcook42/morgoth/detector"
	"github.com/nvcook42/morgoth/engine"
	"github.com/nvcook42/morgoth/metric/types"
	"github.com/nvcook42/morgoth/notifier"
	"github.com/nvcook42/morgoth/schedule"
)

type Supervisor interface {
	GetPattern() Pattern
	AddMetric(types.MetricID)
	Start()
}

type SupervisorStruct struct {
	pattern Pattern
	writer engine.Writer
	detectors []detector.Detector
	notifiers []notifier.Notifier
	schedule schedule.Schedule
}

func NewSupervisor(
	pattern Pattern,
	writer engine.Writer,
	detectors []detector.Detector,
	notifiers []notifier.Notifier,
	schedule schedule.Schedule,
) *SupervisorStruct {

	s := &SupervisorStruct{
		pattern: pattern,
		writer: writer,
		detectors: detectors,
		notifiers: notifiers,
		schedule: schedule,
	}
	return s
}


func (self *SupervisorStruct) GetPattern() Pattern {
	return self.pattern
}

func (self *SupervisorStruct) AddMetric(types.MetricID) {

}

func (self *SupervisorStruct) Start() {

}
