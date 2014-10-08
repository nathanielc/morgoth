package rest_test

import (
	"github.com/nvcook42/morgoth/fitting/rest"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestRestConfShouldDefault(t *testing.T) {
	assert := assert.New(t)

	rc := rest.RESTConf{}

	rc.Default()

	assert.Equal(8000, rc.Port)

}

func TestRestConfShouldValidate(t *testing.T) {
	assert := assert.New(t)

	rc := rest.RESTConf{
		Port: 42,
	}

	err := rc.Validate()
	assert.Nil(err)

	assert.Equal(42, rc.Port)

}

func TestRestConfShouldParse(t *testing.T) {
	assert := assert.New(t)

	var data string = `---
port: 43
`

	rc := rest.RESTConf{}

	err := yaml.Unmarshal([]byte(data), &rc)

	assert.Nil(err)

	assert.Equal(43, rc.Port)

}
