// Define exported types from `metric' to avoid cyclic imports
package types

import (
	"github.com/nathanielc/morgoth/schedule"
)

const (
	MetricPrefix  = "m."
	AnomalyPrefix = "a."
)

type MetricID string

// Return the full metric path to the raw unrotated data
func (self *MetricID) GetRawPath() string {
	return MetricPrefix + string(*self)
}

// Return the full metric path with all appropriate prefixes for a
// a given rotation of the metric data
// If rotation is nil assume no rotation and the raw path is returned
func (self *MetricID) GetRotationPath(rotation *schedule.Rotation) string {
	if rotation != nil {
		return rotation.GetPrefix() + self.GetRawPath()
	}
	return self.GetRawPath()
}

//Return the full metric path for the anomalies of the metric
func (self *MetricID) GetAnomalyPath() string {
	return AnomalyPrefix + string(*self)
}
