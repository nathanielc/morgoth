package influxdb_test

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/nathanielc/morgoth/engine/influxdb"
	"testing"
)

func TestInfluxDBConfShouldDefaultEmpty(t *testing.T) {
	assert := assert.New(t)

	conf := influxdb.InfluxDBConf{}
	conf.Default()

	assert.Equal(conf.Host, "localhost")
	assert.Equal(conf.Port, 8083)

}

func TestInfluxDBConfShouldDefaultNonEmpty(t *testing.T) {
	assert := assert.New(t)

	conf := influxdb.InfluxDBConf{
		Port:     65536,
		Database: "morgoth",
	}
	conf.Default()

	assert.Equal(conf.Host, "localhost")
	assert.Equal(conf.Port, 8083)
	assert.Equal(conf.Database, "morgoth")

}

func TestInfluxDBConfDefaultShouldIgnoreValidFields(t *testing.T) {
	assert := assert.New(t)

	conf := influxdb.InfluxDBConf{
		Host:     "influx",
		Port:     42,
		Database: "morgoth",
	}
	conf.Default()

	assert.Equal(conf.Host, "influx")
	assert.Equal(conf.Port, 42)
	assert.Equal(conf.Database, "morgoth")

}

func TestInfluxDBConfShouldParse(t *testing.T) {
	assert := assert.New(t)

	ic := influxdb.InfluxDBConf{}

	var data string = `---
host: influx1.example.com
port: 4242
user: morgoth
password: mysecret
database: morgothdb
`

	err := yaml.Unmarshal([]byte(data), &ic)

	assert.Nil(err)

	assert.Equal("influx1.example.com", ic.Host)
	assert.Equal(4242, ic.Port)
	assert.Equal("morgoth", ic.User)
	assert.Equal("mysecret", ic.Password)
	assert.Equal("morgothdb", ic.Database)

}
