package rest_test

import (
	log "github.com/cihub/seelog"
	"github.com/nu7hatch/gouuid"
	"github.com/nvcook42/morgoth/engine"
	"github.com/nvcook42/morgoth/fitting"
	_ "github.com/nvcook42/morgoth/fitting/list"
	metric "github.com/nvcook42/morgoth/metric/types"
	app "github.com/nvcook42/morgoth/mocks/app"
	mengine "github.com/nvcook42/morgoth/mocks/engine"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestRestShouldBeRegistered(t *testing.T) {
	assert := assert.New(t)

	_, err := fitting.Registery.GetFactory("rest")

	assert.Nil(err)
}

func TestRestShouldReturnMetrics(t *testing.T) {
	defer log.Flush()
	assert := assert.New(t)

	var yaml = `---
rest:
  port: 7070
`

	app := new(app.App)
	reader := new(mengine.Reader)
	writer := new(mengine.Reader)

	app.On("GetReader").Return(reader).Once()
	app.On("GetWriter").Return(writer).Once()

	reader.On("GetMetrics").
		Return([]metric.MetricID{"m1", "m2"}).
		Once()

	conf, err := fitting.FromYAML(yaml)
	if !assert.Nil(err) {
		assert.Fail(err.Error())
	}
	rest, err := conf.GetFitting()
	if !assert.Nil(err) {
		assert.Fail(err.Error())
	}

	go rest.Start(app)
	time.Sleep(time.Millisecond)

	resp, err := http.Get("http://localhost:7070/metrics")
	if !assert.Nil(err) {
		assert.Fail(err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if !assert.Nil(err) {
		assert.Fail(err.Error())
	}
	assert.Equal(200, resp.StatusCode)
	assert.Equal(
		"{\n  \"metrics\": [\n    \"m1\",\n    \"m2\"\n  ]\n}",
		string(body),
	)
}

func TestRestShouldReturnMetricData(t *testing.T) {
	defer log.Flush()
	assert := assert.New(t)

	var yaml = `---
rest:
  port: 7071
`

	app := new(app.App)
	reader := new(mengine.Reader)
	writer := new(mengine.Reader)

	app.On("GetReader").Return(reader).Once()
	app.On("GetWriter").Return(writer).Once()

	metricID := metric.MetricID("m1")

	zeroTime := time.Time{}
	start := zeroTime.Add(time.Second)
	stop := zeroTime.Add(time.Hour * 2)

	reader.On("GetData", metricID, start, stop, 0).
		Return([]engine.Point{
		{start, 1.0},
		{start.Add(time.Minute), 2.0},
		{start.Add(time.Hour), 3.0},
	}).
		Once()

	conf, err := fitting.FromYAML(yaml)
	if !assert.Nil(err) {
		assert.Fail(err.Error())
	}
	rest, err := conf.GetFitting()
	if !assert.Nil(err) {
		assert.Fail(err.Error())
	}

	go rest.Start(app)
	time.Sleep(time.Millisecond)

	params := url.Values{}
	params.Set("start", start.Format("1405544146"))
	params.Set("stop", stop.Format("1405544146"))
	resp, err := http.Get("http://localhost:7071/data/m1?" + params.Encode())
	if !assert.Nil(err) {
		assert.Fail(err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if !assert.Nil(err) {
		assert.Fail(err.Error())
	}
	assert.Equal(200, resp.StatusCode)
	assert.Equal(
		"{\n  \"data\": [\n    [\n      \"0001-01-01T00:00:00Z\",\n      1\n    ],\n    [\n      \"0001-01-01T00:01:00Z\",\n      2\n    ],\n    [\n      \"0001-01-01T01:00:00Z\",\n      3\n    ]\n  ],\n  \"metric\": \"m1\"\n}",
		string(body),
	)
}

func TestRestShouldReturnMetricAnomalies(t *testing.T) {
	defer log.Flush()
	assert := assert.New(t)

	var yaml = `---
rest:
  port: 7072
`

	app := new(app.App)
	reader := new(mengine.Reader)
	writer := new(mengine.Reader)

	app.On("GetReader").Return(reader).Once()
	app.On("GetWriter").Return(writer).Once()

	metricID := metric.MetricID("m1")

	zeroTime := time.Time{}
	u1, _ := uuid.ParseHex("6ab63420-8857-4f41-429d-5e3ea63944f6")
	u2, _ := uuid.ParseHex("202c4521-ff5a-4038-47be-95d2dc5b79a8")
	u3, _ := uuid.ParseHex("aaf47e42-a458-4e63-4ede-f5d3ea5c649e")
	reader.On("GetAnomalies", metricID, zeroTime, zeroTime.Add(time.Hour)).
		Return([]engine.Anomaly{
		{u1, zeroTime.Add(1 * time.Minute), zeroTime.Add(2 * time.Minute)},
		{u2, zeroTime.Add(3 * time.Minute), zeroTime.Add(4 * time.Minute)},
		{u3, zeroTime.Add(5 * time.Minute), zeroTime.Add(6 * time.Minute)},
	}).
		Once()

	conf, err := fitting.FromYAML(yaml)
	if !assert.Nil(err) {
		assert.Fail(err.Error())
	}
	rest, err := conf.GetFitting()
	if !assert.Nil(err) {
		assert.Fail(err.Error())
	}

	go rest.Start(app)
	time.Sleep(time.Millisecond)

	resp, err := http.Get("http://localhost:7072/anomalies/m1")
	if !assert.Nil(err) {
		assert.Fail(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if !assert.Nil(err) {
		assert.Fail(err.Error())
	}
	assert.Equal(200, resp.StatusCode)
	assert.Equal(
		"{\n  \"data\": [\n    {\n      \"id\": \"6ab63420-8857-4f41-429d-5e3ea63944f6\",\n      \"start\": \"0001-01-01T00:01:00Z\",\n      \"stop\": \"0001-01-01T00:02:00Z\"\n    },\n    {\n      \"id\": \"202c4521-ff5a-4038-47be-95d2dc5b79a8\",\n      \"start\": \"0001-01-01T00:03:00Z\",\n      \"stop\": \"0001-01-01T00:04:00Z\"\n    },\n    {\n      \"id\": \"aaf47e42-a458-4e63-4ede-f5d3ea5c649e\",\n      \"start\": \"0001-01-01T00:05:00Z\",\n      \"stop\": \"0001-01-01T00:06:00Z\"\n    }\n  ],\n  \"metric\": \"m1\"\n}",
		string(body),
	)

}
