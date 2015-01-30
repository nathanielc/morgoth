package app

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nvcook42/morgoth/config"
	mtypes "github.com/nvcook42/morgoth/metric/types"
	engine "github.com/nvcook42/morgoth/mocks/engine"
	metric "github.com/nvcook42/morgoth/mocks/metric"
	"testing"
	"time"
)

func TestAppInit(t *testing.T) {
	assert := assert.New(t)

	var data = `---
#Full App Config example
engine:
  influxdb:
    user: morgoth
    password: morgoth
    database: morgoth

metrics:
  - pattern: .*
    detectors:
    schedule:
      period: 60
      duration: 60

fittings:
  - rest:
      port: 42
`

	config, err := config.Load([]byte(data))
	assert.Nil(err)

	app := New(config)
	assert.NotNil(app)

}

func TestAppShouldNotifyManagerOfNewMetrics(t *testing.T) {
	assert := assert.New(t)

	metricTime := time.Now()
	var metricID mtypes.MetricID = "metric_id"
	metricValue := 42.0

	writer := new(engine.Writer)
	engine := new(engine.Engine)

	manager := new(metric.Manager)

	app := &App{
		engine:  engine,
		manager: manager,
	}

	engine.On("GetWriter").Return(writer).Once()
	writer.On("Insert", metricTime, metricID, metricValue).Return().Once()

	manager.On("NewMetric", metricID).Return().Once()

	appWriter := app.GetWriter()
	if !assert.NotNil(appWriter) {
		assert.Fail("appWriter is nil cannot continue")
	}
	appWriter.Insert(metricTime, metricID, metricValue)

	engine.AssertExpectations(t)
	writer.AssertExpectations(t)
	manager.AssertExpectations(t)

}
