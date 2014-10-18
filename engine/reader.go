package engine

import (
	metric "github.com/nvcook42/morgoth/metric/types"
	"time"
)

type Reader interface {
	GetMetrics() []metric.MetricID
	GetData(metric metric.MetricID, start, stop time.Time, step time.Duration) []Point
	GetAnomalies(metric metric.MetricID, start, stop time.Time) []Anomaly
	GetHistogram(metric metric.MetricID, nbins uint, start, stop time.Time) Histogram
	GetPercentile(metric metric.MetricID, percentile float64, start, stop time.Time) float64
}
