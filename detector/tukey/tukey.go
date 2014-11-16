package tukey

import (
	"fmt"
	"github.com/golang/glog"
	app "github.com/nvcook42/morgoth/app/types"
	"github.com/nvcook42/morgoth/engine"
	metric "github.com/nvcook42/morgoth/metric/types"
	"github.com/nvcook42/morgoth/schedule"
	"sort"
	"time"
)

type Tukey struct {
	rotation  schedule.Rotation
	reader    engine.Reader
	writer    engine.Writer
	threshold float64
}

func (self *Tukey) GetIdentifier() string {
	return fmt.Sprintf(
		"%stukey_%f",
		self.rotation.GetPrefix(),
		self.threshold,
	)
}

func (self *Tukey) Initialize(app app.App, rotation schedule.Rotation) error {
	self.rotation = rotation
	self.reader = app.GetReader()
	self.writer = app.GetWriter()
	return nil
}

func (self *Tukey) Detect(metric metric.MetricID, start, stop time.Time) bool {
	points := self.reader.GetData(&self.rotation, metric, start, stop)
	count := len(points)
	if count == 0 {
		return false
	}

	data := make([]float64, count)
	for i, point := range points {
		data[i] = point.Value
	}
	sort.Float64s(data)

	quater := count / 4
	q1 := data[quater]
	q3 := data[3*quater]

	inner := q3 - q1
	lower := q1 - self.threshold*inner
	upper := q3 + self.threshold*inner

	anomalous := data[0] < lower || data[count-1] > upper

	if glog.V(3) {
		glog.Infof(
			"Metric: %s Count: %d Q1: %f Q3: %f Lower: %f Upper: %f Min: %f Max: %f Anomalous: %v",
			string(metric),
			count,
			q1,
			q3,
			lower,
			upper,
			data[0],
			data[count-1],
			anomalous,
		)
	}

	return anomalous
}
