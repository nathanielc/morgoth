package engine

import (
	metric "github.com/nvcook42/morgoth/metric/types"
	"github.com/nvcook42/morgoth/registery"
	"time"
)

type Engine interface {
	Initialize() error
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
	Metric metric.MetricID
	Start  time.Time
	Stop   time.Time
}

type Histogram struct {
	Bins []float64
	Count int
}
