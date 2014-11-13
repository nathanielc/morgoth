// Define exported types from `metric' to avoid cyclic imports
package types

import (
	"github.com/nvcook42/morgoth/schedule"
	"time"
)

// A metric manager is responsible for creating
// new metric supervisors when a new metric arrives
type Manager interface {
	// Inform the manager of a new metric
	//
	// NOTE: the manager may already know about the 'new' metric
	// in which case it will ignore it.
	NewMetric(MetricID)

	// Instruct manager to perform a detection of all metrics
	// for a given rotation and window
	Detect(rotation schedule.Rotation, start time.Time, stop time.Time)
}
