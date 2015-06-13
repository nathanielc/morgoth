package influxdb_test

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nathanielc/morgoth/engine"
	_ "github.com/nathanielc/morgoth/engine/list"
	"testing"
)

func TestInfluxDBShouldBeRegistered(t *testing.T) {
	assert := assert.New(t)

	_, err := engine.Registery.GetFactory("influxdb")
	assert.Nil(err)
}
