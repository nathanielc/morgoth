package metric

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/metric/set"
	metric "github.com/nvcook42/morgoth/metric/types"
	"github.com/nvcook42/morgoth/schedule"
	"regexp"
	"time"
)

type ManagerStruct struct {
	metrics     *set.Set
	supervisors []pair
	schedule    *schedule.Schedule
}

type pair struct {
	Regexp     *regexp.Regexp
	Supervisor Supervisor
}

func NewManager(schedule *schedule.Schedule, supervisors []Supervisor) metric.Manager {
	m := &ManagerStruct{
		schedule:    schedule,
		metrics:     set.New(0),
		supervisors: make([]pair, len(supervisors)),
	}

	for i := range supervisors {
		supervisor := supervisors[i]
		pattern := supervisor.GetPattern()
		re, err := regexp.Compile(string(pattern))
		if err == nil {
			p := pair{
				re,
				supervisor,
			}
			m.supervisors[i] = p
		} else {
			log.Errorf("Invalid regex for pattern '%s' Error: %s", pattern, err.Error())
		}
	}

	//Start schedule
	m.schedule.Callback = m.detect
	m.schedule.Start()

	return m
}

func (self *ManagerStruct) NewMetric(id metric.MetricID) {
	if !self.metrics.Has(id) {
		supervisor := self.matchSupervisor(id)
		if supervisor == nil {
			log.Warnf("No matching metric pattern for metric '%s'", id)
		}
		supervisor.AddMetric(id)
		self.metrics.Add(id)
	}
}

func (self *ManagerStruct) matchSupervisor(id metric.MetricID) Supervisor {
	for _, pair := range self.supervisors {
		if pair.Regexp.Match([]byte(id)) {
			return pair.Supervisor
		}
	}
	return nil
}

func (self *ManagerStruct) detect(rotation schedule.Rotation, start, stop time.Time) {
	for _, pair := range self.supervisors {
		go pair.Supervisor.Detect(rotation, start, stop)
	}

}
