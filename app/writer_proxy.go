package app

import (
	metric "github.com/nvcook42/morgoth/metric/types"
	"github.com/nvcook42/morgoth/engine"
	"time"
)

type writerProxy struct {
	writer engine.Writer
	manager metric.Manager
}


//
// Proxy Insert method to engine.Writer after informing the manager
//
func (self *writerProxy) Insert(datetime time.Time, metric metric.MetricID, value float64) {
	self.manager.NewMetric(metric)
	self.writer.Insert(datetime, metric, value)
}
