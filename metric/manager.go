package metric

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/metric/set"
	metric "github.com/nvcook42/morgoth/metric/types"
	"regexp"
)

type ManagerStruct struct {
	metrics     *set.Set
	supervisors []pair
}

type pair struct {
	Regexp *regexp.Regexp
	Supervisor Supervisor
}

func NewManager(supervisors []Supervisor) metric.Manager {
	log.Debugf("NewManager: %v", supervisors)
	m := &ManagerStruct{
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

	return m
}

func (self *ManagerStruct) NewMetric(id metric.MetricID) {
	if !self.metrics.Has(id) {
		supervisor := self.matchSupervisor(id)
		if supervisor == nil {
			log.Warnf("No matching metric pattern for metric '%s'", id)
		}
		supervisor.AddMetric(id)
		supervisor.Start()
		self.metrics.Add(id)
	}
}

func (self *ManagerStruct) matchSupervisor(id metric.MetricID) Supervisor {
	for i := range self.supervisors{
		pair := self.supervisors[i]
		if pair.Regexp.Match([]byte(id)) {
			return pair.Supervisor
		}
	}
	return nil
}
