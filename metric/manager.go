package metric

import (
	"github.com/nvcook42/morgoth/metric/set"
	app "github.com/nvcook42/morgoth/app/types"
	mtypes "github.com/nvcook42/morgoth/metric/types"
)

type Manager struct {
	metrics *set.Set
	app app.App
}

func NewManager(metrics []MetricConf, app app.App) *Manager {
	m := new(Manager)
	m.metrics = set.New(0)
	m.app = app
	return m
}

func (self *Manager) NewMetric(id mtypes.MetricID) {
	if !self.metrics.Has(id) {
		//Found real new metric
	}
}

