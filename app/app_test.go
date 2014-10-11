package app_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/nvcook42/morgoth/app"
	"github.com/nvcook42/morgoth/config"
)

func TestAppStart(t *testing.T) {
	assert := assert.New(t)

	var data = `---
#Full App Config example
engine:
  influxdb:
    user: morgoth
    password: morgoth
    database: morgoth

metrics:
  - pattern: .*
    detectors:
    schedule:
      period: 60
      duration: 60

fittings:
  - rest:
      port: 42
`

	config, err := config.Load([]byte(data))
	assert.Nil(err)

	app := app.New(config)
	assert.NotNil(app)

	err = app.Run()
	assert.Nil(err)
	
}


