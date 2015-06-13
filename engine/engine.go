package engine

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/nu7hatch/gouuid"
	"github.com/nathanielc/morgoth/registery"
	"github.com/nathanielc/morgoth/schedule"
	"github.com/nathanielc/morgoth/window"
	"time"
)

type Engine interface {
	Initialize() error
	ConfigureSchedule(schedule *schedule.Schedule) error
	GetReader() Reader
	GetWriter() Writer
	ExecuteQuery(query string) ([]*window.Window, error)
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
