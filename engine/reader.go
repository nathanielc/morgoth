package engine

import (
	metric "github.com/nvcook42/morgoth/metric/types"
	"github.com/nvcook42/morgoth/schedule"
	"time"
)

type Reader interface {
	GetMetrics() []metric.MetricID
	GetData(rotation *schedule.Rotation, metric metric.MetricID, start time.Time, stop time.Time) []Point
	GetAnomalies(metric metric.MetricID, start time.Time, stop time.Time) []Anomaly
	GetHistogram(rotation *schedule.Rotation, metric metric.MetricID, nbins uint, start time.Time, stop time.Time, min float64, max float64) *Histogram
	GetPercentile(rotation *schedule.Rotation, metric metric.MetricID, percentile float64, start time.Time, stop time.Time) float64
	GetDoc(key string) []byte
}
