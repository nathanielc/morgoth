// Define exported types from `metric' to avoid cyclic imports
package types

type MetricID string

// A metric manager is responsible for creating
// new metric supervisors when a new metric arrives
type Manager interface {
	// Inform the manager of a new metric
	//
	// NOTE: the manager may already know about the 'new' metric
	// in which case it will ignore it.
	NewMetric(MetricID)
}
