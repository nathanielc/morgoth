package graphite_test

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/fitting"
	_ "github.com/nvcook42/morgoth/fitting/list"
	metric "github.com/nvcook42/morgoth/metric/types"
	app "github.com/nvcook42/morgoth/mocks/app"
	mengine "github.com/nvcook42/morgoth/mocks/engine"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestRestShouldBeRegistered(t *testing.T) {
	assert := assert.New(t)

	_, err := fitting.Registery.GetFactory("graphite")

	assert.Nil(err)
}

func TestGraphiteShouldWriteData(t *testing.T) {
	defer log.Flush()
	assert := assert.New(t)

	var yaml = `---
graphite:
  port: 2004
`

	app := new(app.App)
	writer := new(mengine.Writer)

	app.On("GetWriter").Return(writer).Once()

	metricID := metric.MetricID("m1")

	zeroTime := time.Time{}
	tm := zeroTime.Add(time.Second).UTC()

	value := 42.0

	writer.On("Insert", tm, metricID, value).Return().Once()

	conf, err := fitting.FromYAML(yaml)
	if !assert.Nil(err) {
		assert.Fail(err.Error())
	}
	graphite, err := conf.GetFitting()
	if !assert.Nil(err) {
		assert.Fail(err.Error())
	}

	go graphite.Start(app)
	time.Sleep(time.Millisecond)

	conn, err := net.Dial("tcp", "localhost:2004")
	if !assert.Nil(err) {
		assert.Fail(err.Error())
	}

	fmt.Fprintf(conn, "%s %f %d\n", metricID, value, tm.Unix())

	conn.Close()
	time.Sleep(time.Millisecond)

	app.AssertExpectations(t)
	writer.AssertExpectations(t)
}
