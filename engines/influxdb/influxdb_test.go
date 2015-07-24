package influxdb_test

import (
	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	_ "github.com/nathanielc/morgoth/engines"
	"testing"
)

func TestInfluxDBShouldBeRegistered(t *testing.T) {
	assert := assert.New(t)

	_, err := morgoth.EngineRegistery.GetFactory("influxdb")
	assert.Nil(err)
}
