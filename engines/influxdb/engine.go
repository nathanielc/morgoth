package influxdb

import (
	"encoding/json"
	"errors"
	"github.com/influxdb/influxdb/client"
	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
)

type InfluxDBEngine struct {
	conf               client.Config
	database           string
	anomalyMeasurement string
	measurementTag     string
}

func (self *InfluxDBEngine) Initialize() error {
	return nil
}

func (self *InfluxDBEngine) NewQueryBuilder(queryTemplate string) (morgoth.QueryBuilder, error) {
	return NewQueryBuilder(queryTemplate)
}

func (self *InfluxDBEngine) RecordAnomalous(w *morgoth.Window) error {
	w.Tags[self.measurementTag] = w.Name
	point := client.Point{
		Measurement: self.anomalyMeasurement,
		Tags:        w.Tags,
		Fields: map[string]interface{}{
			"start": w.Start.Unix(),
			"stop":  w.Stop.Unix(),
		},
		Time:      w.Start,
		Precision: "s",
	}

	batch := client.BatchPoints{
		Points:   []client.Point{point},
		Database: self.database,
	}

	con, err := client.NewClient(self.conf)
	if err != nil {
		return err
	}
	_, err = con.Write(batch)
	return err
}

func (self *InfluxDBEngine) GetWindows(query morgoth.Query) ([]*morgoth.Window, error) {

	con, err := client.NewClient(self.conf)
	if err != nil {
		return nil, err
	}

	q := client.Query{
		Command:  query.Command,
		Database: self.database,
	}

	glog.V(3).Infof("Q: %s", query)
	response, err := con.Query(q)
	if err != nil {
		return nil, err
	}
	if response.Err != nil {
		return nil, response.Err
	}
	windowCount := 0
	for _, result := range response.Results {
		if result.Err != nil {
			return nil, result.Err
		}
		windowCount += len(result.Series)
	}
	windows := make([]*morgoth.Window, windowCount)

	glog.V(3).Infof("Results: %v", response.Results)

	i := 0
	for _, result := range response.Results {
		for _, row := range result.Series {
			w := &morgoth.Window{
				Name:  row.Name,
				Data:  make([]float64, len(row.Values)),
				Tags:  row.Tags,
				Start: query.Start,
				Stop:  query.Stop,
			}

			//Find non time column
			if len(row.Columns) != 2 {
				return nil, errors.New("Queries must select only two columns, a time column and a numeric column")
			}
			numberColumn := 0
			for c, name := range row.Columns {
				if name != "time" {
					numberColumn = c
					break
				}
			}

			for j, point := range row.Values {
				//We care only about the value not the time
				//TODO check columns
				value := point[numberColumn]
				if value != nil {
					p := value.(json.Number)
					v, err := p.Float64()
					if err != nil {
						return nil, err
					}
					w.Data[j] = v
				}
			}

			windows[i] = w
			i++
			glog.V(4).Infof("W: %v", w)
		}
	}

	glog.V(3).Infof("Windows: %v", windows)
	return windows, nil
}
