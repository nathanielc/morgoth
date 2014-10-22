package engine

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/engine"
	_ "github.com/nvcook42/morgoth/engine/list"
	metric "github.com/nvcook42/morgoth/metric/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type EngineTestSuite struct {
	Engine engine.Engine
}

func (self EngineTestSuite) TestAll(t *testing.T) {
	defer log.Flush()

	self.TestGetReader(t)
	self.TestWriteReadOnePoint(t)
	self.TestRecordAnomaly(t)
}

func (self EngineTestSuite) TestGetReader(t *testing.T) {
	assert := assert.New(t)

	reader := self.Engine.GetReader()
	assert.NotNil(reader)
}

func (self EngineTestSuite) TestGetWriter(t *testing.T) {
	assert := assert.New(t)

	writer := self.Engine.GetWriter()
	assert.NotNil(writer)
}

func (self EngineTestSuite) TestWriteReadOnePoint(t *testing.T) {
	assert := assert.New(t)

	metricTime := time.Now().UTC()
	var metricID metric.MetricID = "metric_id"
	metricValue := 42.0

	writer := self.Engine.GetWriter()
	assert.NotNil(writer)

	reader := self.Engine.GetReader()
	assert.NotNil(reader)

	//Delete metric first
	writer.DeleteMetric(metricID)

	//Insert single data piont
	writer.Insert(metricTime, metricID, metricValue)

	metrics := reader.GetMetrics()
	assert.Equal(1, len(metrics))
	assert.Contains(metrics, metricID)

	start := metricTime.Add(-time.Second * 2)
	stop := metricTime.Add(time.Second * 2)
	data := reader.GetData(metricID, start, stop, 0)

	if assert.Equal(1, len(data)) {
		assert.Equal(metricTime, data[0].Time)
		assert.Equal(metricValue, data[0].Value)
	}

}

func (self EngineTestSuite) TestRecordAnomaly(t *testing.T) {

	assert := assert.New(t)

	var metricID metric.MetricID = "metric_id"
	stop := time.Time{}.UTC()
	start := stop.Add(-time.Second * 60)

	writer := self.Engine.GetWriter()
	assert.NotNil(writer)

	reader := self.Engine.GetReader()
	assert.NotNil(reader)

	writer.RecordAnomalous(metricID, start, stop)

	anomalies := reader.GetAnomalies(metricID, start, stop)

	if assert.Equal(1, len(anomalies)) {
		assert.NotNil(anomalies[0].UUID)
		assert.Equal(start, anomalies[0].Start)
		assert.Equal(stop, anomalies[0].Stop)
	}

	newStart := start.Add(time.Hour)
	newStop := stop.Add(time.Hour)

	anomalies = reader.GetAnomalies(metricID, newStart, newStop)
	assert.Equal(0, len(anomalies))
}
