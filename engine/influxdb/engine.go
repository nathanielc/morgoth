package influxdb

import (
	"encoding/json"
	"net/url"
	//"fmt"
	//"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	//"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/influxdb/influxdb/client"
	//"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/nu7hatch/gouuid"
	"github.com/influxdb/influxdb/client"
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nvcook42/morgoth/engine"
	metric "github.com/nvcook42/morgoth/metric/types"
	"github.com/nvcook42/morgoth/schedule"
	"github.com/nvcook42/morgoth/window"
	//"math"
	//"regexp"
	//"strings"
	"time"
)

type InfluxDBEngine struct {
	config *InfluxDBConf
	//client *client.Client
}

func (self *InfluxDBEngine) ExecuteQuery(query string) ([]*window.Window, error) {

	u, err := url.Parse("http://localhost:8086")
	if err != nil {
		return nil, err
	}
	conf := client.Config{
		URL:      *u,
		Username: "root",
		Password: "root",
	}

	con, err := client.NewClient(conf)
	if err != nil {
		return nil, err
	}
	dur, ver, err := con.Ping()
	if err != nil {
		return nil, err
	}
	glog.Infof("Connected version: %s dur: %s", ver, dur)

	q := client.Query{
		Command:  query,
		Database: "mydb",
	}

	glog.Infof("Q: %s", query)
	response, err := con.Query(q)
	if err != nil {
		return nil, err
	}

	glog.Infof("Results: %v", response.Results)

	result := response.Results[0]
	windows := make([]*window.Window, len(result.Series))
	for i, row := range result.Series {
		w := &window.Window{
			Name: row.Name,
			Data: make([]float64, len(row.Values)),
			Tags: row.Tags,
		}

		for j, point := range row.Values {
			//We care only about the value not the time
			//TODO check columns
			p := point[1].(json.Number)
			v, err := p.Float64()
			if err != nil {
				return nil, err
			}
			w.Data[j] = v
		}

		windows[i] = w
		glog.Infof("W: %v", w)
	}

	glog.Infof("Windows: %v", windows)
	return windows, nil
}

func (self *InfluxDBEngine) Initialize() error {
	//client, err := connect(self.config)
	//if err != nil {
	//	return err
	//}
	//self.client = client

	return nil
}

func (self *InfluxDBEngine) GetReader() engine.Reader {
	return self
}

func (self *InfluxDBEngine) GetWriter() engine.Writer {
	return self
}

func (self *InfluxDBEngine) ConfigureSchedule(schedule *schedule.Schedule) error {

	//result, err := self.client.Query("list continuous queries")
	//if err != nil {
	//	return err
	//}

	//existing := make([]string, 0)
	//for _, series := range result {
	//	for _, row := range series.GetPoints() {
	//		existing = append(existing, row[2].(string))
	//	}
	//}

	//pattern := regexp.QuoteMeta(metric.MetricPrefix)
	//for _, rotation := range schedule.Rotations {
	//	resolution := int64(math.Ceil(rotation.Resolution.Seconds()))
	//	q := fmt.Sprintf(
	//		"select first(value) as value from /^%s.*/ group by time(%ds) into %s:series_name",
	//		pattern,
	//		resolution,
	//		rotation.GetPrefix(),
	//	)
	//	found := false
	//	for _, e := range existing {
	//		if e == q {
	//			found = true
	//			break
	//		}
	//	}
	//	if found {
	//		continue
	//	}
	//	glog.Infof("Creating continuous query '%s'", q)
	//	_, err = self.client.Query(q)
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil

}

//////////////////////
// Writer Methods
//////////////////////

func (self *InfluxDBEngine) Insert(datetime time.Time, metric metric.MetricID, value float64) {
	//series := new(client.Series)
	//series.Name = metric.GetRawPath()
	//series.Columns = []string{
	//	"time",
	//	"value",
	//}
	//series.Points = [][]interface{}{
	//	[]interface{}{datetime.Unix(), value},
	//}
	//err := self.client.WriteSeriesWithTimePrecision([]*client.Series{series}, client.Second)
	//if err != nil {
	//	glog.Error(err)
	//}
}

func (self *InfluxDBEngine) RecordAnomalous(metric metric.MetricID, start, stop time.Time) {
	//series := new(client.Series)
	//series.Name = metric.GetAnomalyPath()
	//series.Columns = []string{
	//	"time",
	//	"value",
	//	"uuid",
	//}
	//id, err := uuid.NewV4()
	//if err != nil {
	//	glog.Errorf("Error creating ID for anomaly: %s", err)
	//	return
	//}
	//series.Points = [][]interface{}{
	//	[]interface{}{start.Unix(), stop.Unix(), id.String()},
	//}
	//err = self.client.WriteSeriesWithTimePrecision([]*client.Series{series}, client.Second)
	//if err != nil {
	//	glog.Error(err)
	//}
}

func (self *InfluxDBEngine) DeleteMetric(metric metric.MetricID) {
}

//////////////////////
// Reader Methods
//////////////////////

func (self *InfluxDBEngine) GetMetrics() []metric.MetricID {

	//query := fmt.Sprintf(
	//	"list series /^%s/",
	//	metric.MetricPrefix,
	//)

	//result, err := self.client.Query(query)

	//if err != nil {
	//	glog.Error(err.Error())
	//	return []metric.MetricID{}
	//}
	//if len(result) == 0 {
	//	return []metric.MetricID{}
	//}

	//glog.Infof("List series: %v", result[0].GetPoints())
	//points := result[0].GetPoints()
	//metrics := make([]metric.MetricID, len(points))
	//for i, row := range points {
	//	metrics[i] = metric.MetricID(
	//		strings.Replace(row[1].(string), metric.MetricPrefix, "", 1),
	//	)
	//}
	metrics := make([]metric.MetricID, 0)
	return metrics
}
func (self *InfluxDBEngine) GetData(rotation *schedule.Rotation, metric metric.MetricID, start, stop time.Time) []engine.Point {
	//query := fmt.Sprintf(
	//	"select value from %s",
	//	metric.GetRotationPath(rotation),
	//)
	//if !start.IsZero() {
	//	query += fmt.Sprintf(" where time > %ds", start.Unix())
	//}
	//if !stop.IsZero() {
	//	if start.IsZero() {
	//		query += " where "
	//	} else {
	//		query += " and "
	//	}
	//	query += fmt.Sprintf("time < %ds", stop.Unix())
	//}
	//result, err := self.client.Query(query, client.Second)

	//if err != nil {
	//	glog.Error(err.Error())
	//	return []engine.Point{}
	//}
	//if len(result) == 0 {
	//	return []engine.Point{}
	//}

	//series := result[0]
	//points := series.GetPoints()
	//data := make([]engine.Point, len(points))
	data := make([]engine.Point, 0)
	//for i, row := range points {
	//	sec := int64(row[0].(float64))
	//	data[i].Time = time.Unix(sec, 0)
	//	data[i].Value = row[2].(float64)
	//}
	return data
}

func (self *InfluxDBEngine) GetAnomalies(metric metric.MetricID, start, stop time.Time) []engine.Anomaly {
	//query := fmt.Sprintf(
	//	"select time, value, uuid from %s",
	//	metric.GetAnomalyPath(),
	//)
	//if !start.IsZero() {
	//	query += fmt.Sprintf(" where time > %ds", start.Unix())
	//}
	//if !stop.IsZero() {
	//	if start.IsZero() {
	//		query += " where "
	//	} else {
	//		query += " and "
	//	}
	//	query += fmt.Sprintf("time < %ds", stop.Unix())
	//}
	//result, err := self.client.Query(query, client.Second)

	//if err != nil {
	//	glog.Error(err.Error())
	//	return []engine.Anomaly{}
	//}
	//if len(result) == 0 {
	//	return []engine.Anomaly{}
	//}

	//series := result[0]
	//points := series.GetPoints()
	//data := make([]engine.Anomaly, len(points))
	//for i, row := range points {
	//	start := int64(row[0].(float64))
	//	stop := int64(row[3].(float64))
	//	idStr := row[2].(string)
	//	id, err := uuid.ParseHex(idStr)
	//	if err != nil {
	//		glog.Warningf("Failed to parse UUID %s: %s", idStr, err)
	//		continue
	//	}
	//	data[i].UUID = id
	//	data[i].Start = time.Unix(start, 0)
	//	data[i].Stop = time.Unix(stop, 0)
	//}
	data := make([]engine.Anomaly, 0)
	return data
}
func (self *InfluxDBEngine) GetHistogram(rotation *schedule.Rotation, metric metric.MetricID, nbins uint, start, stop time.Time, min, max float64) *engine.Histogram {
	hist := new(engine.Histogram)
	//hist.Min = min
	//hist.Max = max

	//stepSize := (max - min) / float64(nbins)

	//q := fmt.Sprintf("select count(value), histogram(value, %f, %f, %f) from %s where time > %ds and time < %ds",
	//	stepSize,
	//	min,
	//	max,
	//	metric.GetRotationPath(rotation),
	//	start.Unix(),
	//	stop.Unix(),
	//)

	//result, err := self.client.Query(q)

	//if err != nil {
	//	glog.Error(err.Error())
	//	return hist
	//}
	//if len(result) != 1 {
	//	glog.Error("Invalid results returned for Histogram")
	//	return hist
	//}

	//series := result[0]
	//points := series.GetPoints()
	//hist.Bins = make([]float64, nbins)
	//for _, row := range points {
	//	total := row[1].(float64)
	//	bucketStart := row[2].(float64)
	//	count := row[3].(float64)
	//	i := int((bucketStart - min) / stepSize)
	//	if i == int(nbins) { //Handle last bucket including max value
	//		i--
	//	}
	//	hist.Bins[i] += count / total
	//	hist.Count = uint(total)
	//}
	//if hist.Count == 1 {
	//	glog.V(2).Info("Small hist ", q)
	//}

	return hist

}
func (self *InfluxDBEngine) GetPercentile(rotation *schedule.Rotation, metric metric.MetricID, percentile float64, start, stop time.Time) float64 {
	return 0.0
}
