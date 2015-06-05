package engine

import (
	"time"
)

type Query interface {
	QueryForTimeRange(start, stop time.Time)
}

