package graphite_test

import (
	"fmt"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nathanielc/morgoth/fitting"
	_ "github.com/nathanielc/morgoth/fitting/list"
	metric "github.com/nathanielc/morgoth/metric/types"
	app "github.com/nathanielc/morgoth/mocks/app"
	mengine "github.com/nathanielc/morgoth/mocks/engine"
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
	defer glog.Flush()
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
