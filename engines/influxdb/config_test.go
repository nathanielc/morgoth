package influxdb_test

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/nathanielc/morgoth/engines/influxdb"
	"testing"
)

func TestInfluxDBConfShouldDefaultEmpty(t *testing.T) {
	assert := assert.New(t)

	conf := influxdb.InfluxDBConf{}
	conf.Default()

	assert.Equal("localhost", conf.Host)
	assert.Equal(uint(8083), conf.Port)

}

func TestInfluxDBConfShouldDefaultNonEmpty(t *testing.T) {
	assert := assert.New(t)

	conf := influxdb.InfluxDBConf{
		Port:     65536,
		Database: "morgoth",
	}
	conf.Default()

	assert.Equal("localhost", conf.Host)
	assert.Equal(uint(8083), conf.Port)
	assert.Equal("morgoth", conf.Database)

}

func TestInfluxDBConfDefaultShouldIgnoreValidFields(t *testing.T) {
	assert := assert.New(t)

	conf := influxdb.InfluxDBConf{
		Host:     "influx",
		Port:     42,
		Database: "morgoth",
	}
	conf.Default()

	assert.Equal("influx", conf.Host)
	assert.Equal(uint(42), conf.Port)
	assert.Equal("morgoth", conf.Database)

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
	assert.Equal(uint(4242), ic.Port)
	assert.Equal("morgoth", ic.User)
	assert.Equal("mysecret", ic.Password)
	assert.Equal("morgothdb", ic.Database)

}
