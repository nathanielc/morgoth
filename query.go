package morgoth

import (
	"time"
)

type Query struct {
	Command string
	Start   time.Time
	Stop    time.Time
	tags    map[string]string
}

type QueryBuilder interface {
	GetForTimeRange(start, stop time.Time) Query
}

type AlertQuery struct {
	Query     Query
	Threshold float64
	Notifiers []Notifier
	Message   string
}
