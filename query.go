package morgoth

import (
	"time"
)

type Query string

type QueryConstructor interface {
	GetForTimeRange(start, stop time.Time) Query
}
