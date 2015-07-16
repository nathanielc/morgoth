package mocks

import "github.com/nathanielc/morgoth/engine"
import "github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/mock"

import metric "github.com/nathanielc/morgoth/metric/types"
import "github.com/nathanielc/morgoth/schedule"
import "time"

type Reader struct {
	mock.Mock
}

func (m *Reader) GetMetrics() []metric.MetricID {
	ret := m.Called()

	r0 := ret.Get(0).([]metric.MetricID)

	return r0
}
func (m *Reader) GetData(rotation *schedule.Rotation, metric metric.MetricID, start time.Time, stop time.Time) []engine.Point {
	ret := m.Called(rotation, metric, start, stop)

	r0 := ret.Get(0).([]engine.Point)

	return r0
}
func (m *Reader) GetAnomalies(metric metric.MetricID, start time.Time, stop time.Time) []engine.Anomaly {
	ret := m.Called(metric, start, stop)

	r0 := ret.Get(0).([]engine.Anomaly)

	return r0
}
func (m *Reader) GetHistogram(rotation *schedule.Rotation, metric metric.MetricID, nbins uint, start time.Time, stop time.Time, min float64, max float64) *engine.Histogram {
	ret := m.Called(rotation, metric, nbins, start, stop, min, max)

	r0 := ret.Get(0).(*engine.Histogram)

	return r0
}
func (m *Reader) GetPercentile(rotation *schedule.Rotation, metric metric.MetricID, percentile float64, start time.Time, stop time.Time) float64 {
	ret := m.Called(rotation, metric, percentile, start, stop)

	r0 := ret.Get(0).(float64)

	return r0
}
