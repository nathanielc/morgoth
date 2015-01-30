package mgof_test

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/nvcook42/morgoth/detector/mgof"
	"testing"
)

func TestMGOFConfShouldDefault(t *testing.T) {
	assert := assert.New(t)

	mc := mgof.MGOFConf{}

	mc.Default()

	assert.Equal(3, mc.NullConfidence)
	assert.Equal(15, mc.NBins)
	assert.Equal(3, mc.NormalCount)
	assert.Equal(20, mc.MaxFingerprints)
	assert.Equal(0.0, mc.Min)
	assert.Equal(0.0, mc.Max)
}

func TestMGOFConfValidateShouldFailBadMinMax(t *testing.T) {
	assert := assert.New(t)

	mc := mgof.MGOFConf{Min: 0, Max: -1}

	err := mc.Validate()
	assert.NotNil(err)
}

func TestMGOFConfValidateShouldFailBadNullConfidence(t *testing.T) {
	assert := assert.New(t)

	mc := mgof.MGOFConf{NullConfidence: 0}

	err := mc.Validate()
	assert.NotNil(err)
}

func TestMGOFConfValidateShouldFailBadNBins(t *testing.T) {
	assert := assert.New(t)

	mc := mgof.MGOFConf{NBins: 0}

	err := mc.Validate()
	assert.NotNil(err)
}

func TestMGOFConfValidateShouldFailBadNormalCount(t *testing.T) {
	assert := assert.New(t)

	mc := mgof.MGOFConf{NormalCount: 0}

	err := mc.Validate()
	assert.NotNil(err)
}

func TestMGOFConfValidateShouldFailBadMaxFingerprints(t *testing.T) {
	assert := assert.New(t)

	mc := mgof.MGOFConf{MaxFingerprints: 0}

	err := mc.Validate()
	assert.NotNil(err)
}

func TestMGOFConfValidateShouldPass(t *testing.T) {
	assert := assert.New(t)

	mc := mgof.MGOFConf{
		Min:             -1.0,
		Max:             1.0,
		NullConfidence:  2,
		NBins:           5,
		NormalCount:     4,
		MaxFingerprints: 15,
	}

	err := mc.Validate()
	assert.Nil(err)
}

func TestMGOFConfShouldParse(t *testing.T) {
	assert := assert.New(t)

	var data string = `---
min: -10
max: 10
null_confidence: 5
nbins: 5
normal_count: 2
max_fingerprints: 25
`

	mc := mgof.MGOFConf{}

	err := yaml.Unmarshal([]byte(data), &mc)

	assert.Nil(err)

	assert.Equal(-10.0, mc.Min)
	assert.Equal(10.0, mc.Max)
	assert.Equal(5, mc.NullConfidence)
	assert.Equal(5, mc.NBins)
	assert.Equal(2, mc.NormalCount)
	assert.Equal(25, mc.MaxFingerprints)
}
