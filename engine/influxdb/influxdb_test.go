package influxdb_test

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nvcook42/morgoth/engine"
	_ "github.com/nvcook42/morgoth/engine/list"
	"testing"
)

func TestInfluxDBShouldBeRegistered(t *testing.T) {
	assert := assert.New(t)

	_, err := engine.Registery.GetFactory("influxdb")
	assert.Nil(err)
}
