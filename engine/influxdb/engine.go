package influxdb

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/influxdb/influxdb/client"
	"github.com/nvcook42/morgoth/engine"
	metric "github.com/nvcook42/morgoth/metric/types"
	"github.com/nvcook42/morgoth/schedule"
	"math"
	"regexp"
	"time"
)


type InfluxDBEngine struct {
	config *InfluxDBConf
	client *client.Client
}

func (self *InfluxDBEngine) Initialize() error {
	client, err := connect(self.config)
	if err != nil {
		return err
	}
	self.client = client

	return nil
}

func (self *InfluxDBEngine) GetReader() engine.Reader {
	return self
}

func (self *InfluxDBEngine) GetWriter() engine.Writer {
	return self
}

func (self *InfluxDBEngine) ConfigureSchedule(schedule *schedule.Schedule) error {

	result, err := self.client.Query("list continuous queries")
	if err != nil {
		return err
	}

	existing := make([]string, 0)
	for _, series := range result {
		for _, row := range series.GetPoints() {
			existing = append(existing, row[2].(string))
		}
	}

	pattern := regexp.QuoteMeta(metric.MetricPrefix)
	for _, rotation := range schedule.Rotations {
		resolution := int64(math.Ceil(rotation.Resolution.Seconds()))
		q := fmt.Sprintf(
			"select first(value) as value from /^%s.*/ group by time(%ds) into %s:series_name",
			pattern,
			resolution,
			rotation.GetPrefix(),
		)
		found := false
		for _, e := range existing {
			if e == q {
				found = true
				break
			}
		}
		if found {
			continue
		}
		glog.Infof("Creating continuous query '%s'", q)
		_, err = self.client.Query(q)
		if err != nil {
			return err
		}
	}

	return nil

}

//////////////////////
// Writer Methods
//////////////////////

func (self *InfluxDBEngine) Insert(datetime time.Time, metric metric.MetricID, value float64) {
	series := new(client.Series)
	series.Name = metric.GetRawPath()
	series.Columns = []string{
		"time",
		"value",
	}
	series.Points = [][]interface{}{
		[]interface{}{datetime.Unix(), value},
	}
	err := self.client.WriteSeriesWithTimePrecision([]*client.Series{series}, client.Second)
	if err != nil {
		glog.Error(err)
	}
}

func (self *InfluxDBEngine) RecordAnomalous(metric metric.MetricID, start, stop time.Time) {
	series := new(client.Series)
	series.Name = metric.GetAnomalyPath()
	series.Columns = []string{
		"time",
		"value",
	}
	series.Points = [][]interface{}{
		[]interface{}{start.Unix(), float64(stop.Sub(start))},
	}
	err := self.client.WriteSeriesWithTimePrecision([]*client.Series{series}, client.Second)
	if err != nil {
		glog.Error(err)
	}
}

func (self *InfluxDBEngine) DeleteMetric(metric metric.MetricID) {
}


//////////////////////
// Reader Methods
//////////////////////

func (self *InfluxDBEngine) GetMetrics() []metric.MetricID {
	return nil
}
func (self *InfluxDBEngine) GetData(rotation *schedule.Rotation, metric metric.MetricID, start, stop time.Time) []engine.Point {
	result, err := self.client.Query(
		fmt.Sprintf("select value from %s where time > %ds and time < %ds",
			metric.GetRotationPath(rotation),
			start.Unix(),
			stop.Unix(),
		),
		client.Second,
	)

	if err != nil {
		glog.Error(err.Error())
		return []engine.Point{}
	}
	if len(result) == 0 {
		return []engine.Point{}
	}

	series := result[0]
	points := series.GetPoints()
	data := make([]engine.Point, len(points))
	for i, row := range points {
		sec := int64(row[0].(float64))
		data[i].Time = time.Unix(sec, 0)
		data[i].Value = row[2].(float64)
	}
	return data
}

func (self *InfluxDBEngine) GetAnomalies(metric metric.MetricID, start, stop time.Time) []engine.Anomaly {
	return nil
}
func (self *InfluxDBEngine) GetHistogram(rotation *schedule.Rotation, metric metric.MetricID, nbins uint, start, stop time.Time, min, max float64) *engine.Histogram {
	hist := new(engine.Histogram)

	stepSize := (max - min) / float64(nbins)

	q := fmt.Sprintf("select count(value), histogram(value, %f, %f, %f) from %s where time > %ds and time < %ds",
		stepSize,
		min,
		max,
		metric.GetRotationPath(rotation),
		start.Unix(),
		stop.Unix(),
	)

	result, err := self.client.Query(q)

	if err != nil {
		glog.Error(err.Error())
		return hist
	}
	if len(result) != 1 {
		glog.Error("Invalid results returned for Histogram")
		return hist
	}

	series := result[0]
	points := series.GetPoints()
	hist.Bins = make([]float64, nbins)
	for _, row := range points {
		total := row[1].(float64)
		bucketStart := row[2].(float64)
		count := row[3].(float64)
		i := int((bucketStart - min) / stepSize)
		if i == int(nbins) { //Handle last bucket including max value
			i--
		}
		hist.Bins[i] += count / total
		hist.Count = uint(total)
	}
	if hist.Count == 1 {
		glog.V(2).Info("Small hist ", q)
	}

	return hist

}
func (self *InfluxDBEngine) GetPercentile(rotation *schedule.Rotation, metric metric.MetricID, percentile float64, start, stop time.Time) float64 {
	return 0.0
}

