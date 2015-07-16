package tukey_test

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/nathanielc/morgoth/detector/tukey"
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
