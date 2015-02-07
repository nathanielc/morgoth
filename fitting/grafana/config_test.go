package grafana_test

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/nvcook42/morgoth/fitting/grafana"
	_ "github.com/nvcook42/morgoth/engine/list"
	"testing"
)

func TestGrafanaConfShouldDefault(t *testing.T) {
	defer glog.Flush()
	assert := assert.New(t)

	rc := grafana.GrafanaConf{}

	rc.Default()

	assert.Equal("http://grafanarel.s3.amazonaws.com/grafana-1.9.1.tar.gz", rc.URL)
	assert.Equal(8080, rc.Port)
	assert.Equal("grafana_tmp", rc.Dir)
	assert.Equal("grafana", rc.GrafanaDB)
}

func TestGrafanaConfShouldValidate(t *testing.T) {
	defer glog.Flush()
	assert := assert.New(t)

	rc := grafana.GrafanaConf{
		URL: "http://localhost/grafana.tar.gz",
		Port: 42,
		Dir: "/tmp/grafana",
		GrafanaDB: "grafanadb",
	}

	err := rc.Validate()
	assert.Nil(err)

	assert.Equal("http://localhost/grafana.tar.gz", rc.URL)
	assert.Equal(42, rc.Port)
	assert.Equal("/tmp/grafana", rc.Dir)
	assert.Equal("grafanadb", rc.GrafanaDB)

}

func TestGrafanaConfShouldParse(t *testing.T) {
	defer glog.Flush()
	assert := assert.New(t)

	var data string = `---
url: http://localhost/grafana.tar.gz
port: 43
dir: /tmp/grafana
grafana_db: grafanadb
influxdb: &INFLUXDB
  host: localhost
  port: 8086
  user: test_user
  password: test_password
  database: morgoth
`

	rc := grafana.GrafanaConf{}

	err := yaml.Unmarshal([]byte(data), &rc)

	assert.Nil(err)

	assert.Equal("http://localhost/grafana.tar.gz", rc.URL)
	assert.Equal(43, rc.Port)
	assert.Equal("/tmp/grafana", rc.Dir)
	assert.Equal("grafanadb", rc.GrafanaDB)

	assert.Equal("localhost", rc.InfluxDBConf.Host)
	assert.Equal(8086, rc.InfluxDBConf.Port)
	assert.Equal("test_user", rc.InfluxDBConf.User)
	assert.Equal("test_password", rc.InfluxDBConf.Password)
	assert.Equal("morgoth", rc.InfluxDBConf.Database)

}
