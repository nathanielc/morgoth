package notifier

import (
	metric "github.com/nvcook42/morgoth/metric/types"
	"github.com/nvcook42/morgoth/registery"
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
