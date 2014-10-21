package mocks

import "github.com/nvcook42/morgoth/engine"
import "github.com/stretchr/testify/mock"

import metric "github.com/nvcook42/morgoth/metric/types"
import "time"

type Reader struct {
 mock.Mock
}

func (m *Reader) GetMetrics() []metric.MetricID {
 ret := m.Called()

 r0 := ret.Get(0).([]metric.MetricID)

 return r0
}
func (m *Reader) GetData(metric metric.MetricID, start, stop time.Time, step time.Duration) []engine.Point {
 ret := m.Called(metric, start, stop, step)

 r0 := ret.Get(0).([]engine.Point)

 return r0
}
func (m *Reader) GetAnomalies(metric metric.MetricID, start, stop time.Time) []engine.Anomaly {
 ret := m.Called(metric, start, stop)

 r0 := ret.Get(0).([]engine.Anomaly)

 return r0
}
func (m *Reader) GetHistogram(metric metric.MetricID, nbins uint, start, stop time.Time) engine.Histogram {
 ret := m.Called(metric, nbins, start, stop)

 r0 := ret.Get(0).(engine.Histogram)

 return r0
}
func (m *Reader) GetPercentile(metric metric.MetricID, percentile float64, start, stop time.Time) float64 {
 ret := m.Called(metric, percentile, start, stop)

 r0 := ret.Get(0).(float64)

 return r0
}
