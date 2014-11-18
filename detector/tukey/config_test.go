package tukey_test

import (
	"github.com/nvcook42/morgoth/detector/tukey"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestTukeyConfShouldDefault(t *testing.T) {
	assert := assert.New(t)

	mc := tukey.TukeyConf{}

	mc.Default()

	assert.Equal(3.0, mc.Threshold)
}

func TestTukeyConfValidateShouldFailBadThreshold(t *testing.T) {
	assert := assert.New(t)

	mc := tukey.TukeyConf{Threshold: -1.0}

	err := mc.Validate()
	assert.NotNil(err)
}

func TestTukeyConfValidateShouldPass(t *testing.T) {
	assert := assert.New(t)

	mc := tukey.TukeyConf{
		Threshold: 5.0,
	}

	err := mc.Validate()
	assert.Nil(err)
}

func TestTukeyConfShouldParse(t *testing.T) {
	assert := assert.New(t)

	var data string = `---
threshold: 3.5
`

	mc := tukey.TukeyConf{}

	err := yaml.Unmarshal([]byte(data), &mc)

	assert.Nil(err)

	assert.Equal(3.5, mc.Threshold)
}
