package morgoth

import (
	"time"
)

type Query struct {
	Command string
	Start   time.Time
	Stop    time.Time
}

type QueryBuilder interface {
	GetForTimeRange(start, stop time.Time) Query
}
