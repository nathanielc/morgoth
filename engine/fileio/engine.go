package fileio

import (
	"bufio"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/nu7hatch/gouuid"
	"github.com/nvcook42/morgoth/engine"
	metric "github.com/nvcook42/morgoth/metric/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	anomalyPath = "anomalies"
)

type FileIOEngine struct {
	Dir string
}

func (self *FileIOEngine) Initialize() error {

	err := os.MkdirAll(self.Dir, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}

	anomalyDir := filepath.Join(self.Dir, anomalyPath)
	err = os.MkdirAll(anomalyDir, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (self *FileIOEngine) GetReader() engine.Reader {
	return self
}

func (self *FileIOEngine) GetWriter() engine.Writer {
	return self
}

///////////////////////////////
// Writer Methods
///////////////////////////////

func (self *FileIOEngine) Insert(datetime time.Time, metric metric.MetricID, value float64) {
	path := self.pathForMetric(metric)
	file, err := openForAppend(path)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Fprintf(file, "%s %f\n", dateString(datetime), value)
	file.Close()
}

func (self *FileIOEngine) RecordAnomalous(metric metric.MetricID, start, stop time.Time) {
	path := self.pathForMetricAnomalies(metric)
	file, err := openForAppend(path)
	if err != nil {
		log.Error(err)
		return
	}
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Fprintf(file, "%s %f %f\n", uuid.String(), dateString(start), dateString(stop))
	file.Close()
}

func (self *FileIOEngine) DeleteMetric(metric metric.MetricID) {
	path := self.pathForMetric(metric)
	os.Remove(path)
}

///////////////////////////////
// Reader Methods
///////////////////////////////

func (self *FileIOEngine) GetMetrics() []metric.MetricID {
	metrics := make([]metric.MetricID, 0, 10)
	dirInfo, err := ioutil.ReadDir(self.Dir)
	if err != nil {
		return metrics
	}
	for _, info := range dirInfo {
		if !info.IsDir() {
			metrics = append(metrics, metric.MetricID(info.Name()))
		}
	}
	return metrics
}

func (self *FileIOEngine) GetData(metric metric.MetricID, start, stop time.Time, step time.Duration) []engine.Point {
	data := make([]engine.Point, 0, 10)
	path := self.pathForMetric(metric)
	file, err := openForRead(path)
	if err != nil {
		return data
	}
	for file.Scan() {
		line := file.Text()
		parts := strings.Split(line, " ")
		datetime, err := parseDate(parts[0])
		if err != nil {
			log.Error(err)
			continue
		}

		if (datetime.Before(stop) && datetime.After(start)) || // Double bounded
			(start.IsZero() && datetime.Before(stop)) || // Single bound
			(stop.IsZero() && datetime.After(start)) || // Single bound
			(start.IsZero() && stop.IsZero()) { // No bounds

			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				log.Error(err)
				continue
			}
			data = append(data, engine.Point{datetime, value})
		}
	}
	return data
}

func (self *FileIOEngine) GetAnomalies(metric metric.MetricID, start, stop time.Time) []engine.Anomaly {
	data := make([]engine.Anomaly, 0, 10)
	path := self.pathForMetric(metric)
	file, err := openForRead(path)
	if err != nil {
		return data
	}
	for file.Scan() {
		line := file.Text()
		parts := strings.Split(line, " ")
		uuid, err := uuid.ParseHex(parts[0])
		if err != nil {
			log.Error(err)
			continue
		}
		anomalyStart, err := parseDate(parts[1])
		if err != nil {
			log.Error(err)
			continue
		}
		anomalyStop, err := parseDate(parts[2])
		if err != nil {
			log.Error(err)
			continue
		}

		if (anomalyStop.After(start) || anomalyStop.Equal(start)) &&
			(anomalyStart.Before(stop) || anomalyStart.Equal(stop)) {
			data = append(data, engine.Anomaly{uuid, start, stop})
		}
	}
	return data
}

func (self *FileIOEngine) GetHistogram(metric metric.MetricID, nbins uint, start, stop time.Time) engine.Histogram {
	return engine.Histogram{}
}

func (self *FileIOEngine) GetPercentile(metric metric.MetricID, percentile float64, start, stop time.Time) float64 {
	return 0.0
}

///////////////////////////////
// Helper Methods
///////////////////////////////

func (self *FileIOEngine) pathForMetric(metric metric.MetricID) string {
	return filepath.Join(self.Dir, string(metric))
}

func (self *FileIOEngine) pathForMetricAnomalies(metric metric.MetricID) string {
	return filepath.Join(self.Dir, anomalyPath, string(metric))
}

func openForAppend(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
}

func openForRead(path string) (*bufio.Scanner, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return bufio.NewScanner(file), nil
}

func dateString(datetime time.Time) string {
	return datetime.Format(time.RFC3339Nano)
}

func parseDate(datetime string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, datetime)
}
