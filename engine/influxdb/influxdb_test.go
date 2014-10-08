package influxdb_test

import (
	"github.com/nvcook42/morgoth/engine"
	_ "github.com/nvcook42/morgoth/engine/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInfluxDBShouldBeRegistered(t *testing.T) {
	assert := assert.New(t)

	_, err := engine.Registery.GetFactory("influxdb")
	assert.Nil(err)
}
