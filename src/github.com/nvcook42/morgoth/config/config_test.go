package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigShouldNotParseInvalidYAML(t *testing.T) {
	assert := assert.New(t)
	var data = `
---
data_engine:
  mongodb:
	use_sharding: false
`
	config, err := Load([]byte(data))
	assert.NotNil(err)
	assert.Nil(config)
}

func TestConfigShouldNotValidateBadEngineConf(t *testing.T) {
	assert := assert.New(t)

	var data = `
---
data_engine:
  bad_key: 1
`
	_, err := Load([]byte(data))
	assert.NotNil(err)

}

func TestConfigShouldNotValidateBadMetricConf(t *testing.T) {
	assert := assert.New(t)

	var data = `
---
metrics:
  - {}
  - {}
`
	_, err := Load([]byte(data))
	assert.NotNil(err)

}

func TestConfigShouldNotValidateBadFittingConf(t *testing.T) {
	assert := assert.New(t)

	var data = `
---
fittings:
	unknown: {}
	bad: {}
`
	_, err := Load([]byte(data))
	assert.NotNil(err)

}

