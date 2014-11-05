package app

import (
	"github.com/nvcook42/morgoth/engine"
	metric "github.com/nvcook42/morgoth/metric/types"
	"time"
)

type writerProxy struct {
	writer  engine.Writer
	manager metric.Manager
}

//
// Proxy Insert method to engine.Writer after informing the manager
//
func (self *writerProxy) Insert(datetime time.Time, metric metric.MetricID, value float64) {
	self.manager.NewMetric(metric)
	self.writer.Insert(datetime, metric, value)
}

func (self *writerProxy) RecordAnomalous(metric metric.MetricID, start, stop time.Time) {
	self.writer.RecordAnomalous(metric, start, stop)
}

func (self *writerProxy) DeleteMetric(metric metric.MetricID) {
	self.writer.DeleteMetric(metric)
}

