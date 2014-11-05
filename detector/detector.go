package detector

import (
	app "github.com/nvcook42/morgoth/app/types"
	metric "github.com/nvcook42/morgoth/metric/types"
	"github.com/nvcook42/morgoth/registery"
	"github.com/nvcook42/morgoth/schedule"
	"time"
)

type Detector interface {
	Initialize(app.App, schedule.Rotation) error
	Detect(metric metric.MetricID, start, stop time.Time) bool
	GetIdentifier() string
}

var (
	Registery *registery.Registery
)

func init() {
	Registery = registery.New()
}
