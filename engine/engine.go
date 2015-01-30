package engine

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/nu7hatch/gouuid"
	"github.com/nvcook42/morgoth/registery"
	"github.com/nvcook42/morgoth/schedule"
	"time"
)

type Engine interface {
	Initialize() error
	ConfigureSchedule(schedule *schedule.Schedule) error
	GetReader() Reader
	GetWriter() Writer
}

var (
	Registery *registery.Registery
)

func init() {
	Registery = registery.New()
}

type Point struct {
	Time  time.Time
	Value float64
}

type Anomaly struct {
	UUID  *uuid.UUID
	Start time.Time
	Stop  time.Time
}

type Histogram struct {
	Bins  []float64
	Count uint
	Min   float64
	Max   float64
}
