package mocks

import "github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/mock"

import metric "github.com/nvcook42/morgoth/metric/types"
import "time"

type Writer struct {
	mock.Mock
}

func (m *Writer) Insert(datetime time.Time, metric metric.MetricID, value float64) {
	m.Called(datetime, metric, value)
}
func (m *Writer) RecordAnomalous(metric metric.MetricID, start, stop time.Time) {
	m.Called(metric, start)
}
func (m *Writer) DeleteMetric(metric metric.MetricID) {
	m.Called(metric)
}
