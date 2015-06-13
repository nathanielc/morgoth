package notifier

import (
	metric "github.com/nathanielc/morgoth/metric/types"
	"github.com/nathanielc/morgoth/registery"
	"time"
)

type Notifier interface {
	Notify(detectorName string, metric metric.MetricID, start time.Time, stop time.Time)
}

var (
	Registery *registery.Registery
)

func init() {
	Registery = registery.New()
}
