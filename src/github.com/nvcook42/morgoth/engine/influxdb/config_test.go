package influxdb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInfluxDBConfShouldDefaultEmpty(t *testing.T) {
	assert := assert.New(t)

	conf := InfluxDBConf{}
	conf.Default()

	assert.Equal(conf.Host, "localhost")
	assert.Equal(conf.Port, 8083)

}

func TestInfluxDBConfShouldDefaultNonEmpty(t *testing.T) {
	assert := assert.New(t)

	conf := InfluxDBConf{
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

	conf := InfluxDBConf{
		Host:     "influx",
		Port:     42,
		Database: "morgoth",
	}
	conf.Default()

	assert.Equal(conf.Host, "influx")
	assert.Equal(conf.Port, 42)
	assert.Equal(conf.Database, "morgoth")

}
