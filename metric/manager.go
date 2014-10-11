package metric

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/metric/set"
	app "github.com/nvcook42/morgoth/app/types"
	metric "github.com/nvcook42/morgoth/metric/types"
)

type ManagerStruct struct {
	metrics *set.Set
	supervisors map[Pattern]Supervisor
	app app.App
}

func NewManager(supervisors []Supervisor, app app.App) metric.Manager {
	log.Debugf("NewManager: %v, %v", supervisors, app)
	m := &ManagerStruct{
		metrics: set.New(0),
		supervisors: make(map[Pattern]Supervisor, len(supervisors)),
		app: app,
	}

	for i := range supervisors {
		supervisor := supervisors[i]
		pattern := supervisor.GetPattern()
		m.supervisors[pattern] = supervisor
	}

	return m
}

func (self *ManagerStruct) NewMetric(id metric.MetricID) {
	if !self.metrics.Has(id) {
		//Found real new metric
	}
}

