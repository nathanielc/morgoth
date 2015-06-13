package detector

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	app "github.com/nathanielc/morgoth/app/types"
	metric "github.com/nathanielc/morgoth/metric/types"
	"github.com/nathanielc/morgoth/registery"
	"github.com/nathanielc/morgoth/schedule"
	"time"
)

const TraceLevel glog.Level = 3

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
