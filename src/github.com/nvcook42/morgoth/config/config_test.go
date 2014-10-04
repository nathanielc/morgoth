package config

import (
	log "github.com/cihub/seelog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestConfigShouldNotValidateBadDataEngineConf(t *testing.T) {
	assert := assert.New(t)

	var data = `
---
data_engine:
  bad_key: 1
`
	_, err := Load([]byte(data))
	assert.NotNil(err)

}

func TestConfigShouldParseDetectorConf(t *testing.T) {
	assert := assert.New(t)

	var data = `
---
data_engine:
  mongodb:
    database: morgoth
    use_sharding: false
metrics:
  - pattern: .*
    detectors:
      - type: mgof
        mgof:
          chi: 0.08
`

	config, err := Load([]byte(data))
	require.Nil(t, err, "Error loading config: %v\n", err)

	log.Debugf("Config: %v", config)
	assert.Equal(config.Metrics[0].Detectors[0].Type, MGOF)
	assert.Equal(config.Metrics[0].Detectors[0].MGOF.CHI, 0.08)

	assert.Nil(config.Validate())
}

func TestConfigShouldParseScheduleConf(t *testing.T) {
	assert := assert.New(t)

	var data = `
---
data_engine:
  mongodb:
    database: morgoth
    use_sharding: false
metrics:
  - pattern: .*
    schedule:
      duration: 30
      period: 30
      delay: 30
`

	config, err := Load([]byte(data))
	require.Nil(t, err, "Error loading config: %v\n", err)

	log.Debugf("Config: %v", config)
	assert.Equal(config.Metrics[0].Schedule.Duration, 30)
	assert.Equal(config.Metrics[0].Schedule.Period, 30)
	assert.Equal(config.Metrics[0].Schedule.Delay, 30)

	assert.Nil(config.Validate())
}
