package engine

import (
	"github.com/nvcook42/morgoth/registery"
	"github.com/nvcook42/morgoth/schedule"
	"time"
	"github.com/nu7hatch/gouuid"
)

type Engine interface {
	Initialize() error
	ConfigureSchedule(schedule schedule.Schedule) error
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
	UUID   *uuid.UUID
	Start  time.Time
	Stop   time.Time
}

type Histogram struct {
	Bins []float64
	Count uint
}
