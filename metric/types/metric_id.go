// Define exported types from `metric' to avoid cyclic imports
package types

import (
	"github.com/nvcook42/morgoth/schedule"
)

const (
	MetricPrefix  = "m."
	AnomalyPrefix = "a."
)

type MetricID string

func (self *MetricID) GetRawPath() string {
	return MetricPrefix + string(*self)
}

// Return the full metric path with all appropriate prefixes for a
// a given rotation of the metric data
// If rotation is nil assume no rotation
func (self *MetricID) GetMetricPath(rotation *schedule.Rotation) string {
	if rotation != nil {
		return rotation.GetPrefix() + string(*self)
	}

	return self.GetRawPath()
}

//Return the full metric path for the anomalies of the metric
func (self *MetricID) GetAnomalyPath() string {
	return AnomalyPrefix + string(*self)
}
