package mgof_test

import (
	"github.com/nvcook42/morgoth/detector/mgof"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestMGOFConfShouldDefault(t *testing.T) {
	assert := assert.New(t)

	mc := mgof.MGOFConf{}

	mc.Default()

	assert.Equal(0.5, mc.CHI)

}

func TestMGOFConfValidateShouldFailBadCHI(t *testing.T) {
	assert := assert.New(t)

	mc := mgof.MGOFConf{CHI: 0}

	err := mc.Validate()
	assert.NotNil(err)
}

func TestMGOFConfValidateShouldPass(t *testing.T) {
	assert := assert.New(t)

	mc := mgof.MGOFConf{CHI: 0.9}

	err := mc.Validate()
	assert.Nil(err)
}

func TestMGOFConfShouldParse(t *testing.T) {
	assert := assert.New(t)

	var data string = `---
chi: 0.42
`

	mc := mgof.MGOFConf{CHI: 0.9}

	err := yaml.Unmarshal([]byte(data), &mc)

	assert.Nil(err)

	assert.Equal(0.42, mc.CHI)
}
