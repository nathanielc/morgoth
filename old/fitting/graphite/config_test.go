package graphite_test

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/nathanielc/morgoth/fitting/graphite"
	"testing"
)

func TestGraphiteConfShouldDefault(t *testing.T) {
	assert := assert.New(t)

	rc := graphite.GraphiteConf{}

	rc.Default()

	assert.Equal(2003, rc.Port)

}

func TestGraphiteConfShouldValidate(t *testing.T) {
	assert := assert.New(t)

	rc := graphite.GraphiteConf{
		Port: 42,
	}

	err := rc.Validate()
	assert.Nil(err)

	assert.Equal(42, rc.Port)

}

func TestGraphiteConfShouldParse(t *testing.T) {
	assert := assert.New(t)

	var data string = `---
port: 43
`

	rc := graphite.GraphiteConf{}

	err := yaml.Unmarshal([]byte(data), &rc)

	assert.Nil(err)

	assert.Equal(43, rc.Port)

}
