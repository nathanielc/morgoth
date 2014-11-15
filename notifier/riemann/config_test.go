package riemann_test

import (
	"github.com/nvcook42/morgoth/notifier/riemann"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestRiemannConfShouldDefault(t *testing.T) {
	assert := assert.New(t)

	rc := riemann.RiemannConf{}

	rc.Default()

	assert.Equal("localhost", rc.Host)
	assert.Equal(5555, rc.Port)

}

func TestRiemannConfShouldValidate(t *testing.T) {
	assert := assert.New(t)

	rc := riemann.RiemannConf{
		Host: "example.com",
		Port: 42,
	}

	err := rc.Validate()
	assert.Nil(err)

	assert.Equal("example.com", rc.Host)
	assert.Equal(42, rc.Port)

}

func TestRiemannConfShouldParse(t *testing.T) {
	assert := assert.New(t)

	var data string = `---
host: riemann
port: 43
`

	rc := riemann.RiemannConf{}

	err := yaml.Unmarshal([]byte(data), &rc)

	assert.Nil(err)

	assert.Equal("riemann", rc.Host)
	assert.Equal(43, rc.Port)

}
