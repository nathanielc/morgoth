package kstest_test

import (
	"github.com/nvcook42/morgoth/detector/kstest"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestKSTestConfShouldDefault(t *testing.T) {
	assert := assert.New(t)

	ks := kstest.KSTestConf{}

	ks.Default()

	assert.Equal(1, ks.Strictness)
	assert.Equal(3, ks.NormalCount)
	assert.Equal(20, ks.MaxFingerprints)

}

func TestKSTestConfValidateShouldFailBadStrictness(t *testing.T) {
	assert := assert.New(t)

	ks := kstest.KSTestConf{Strictness: 6}

	err := ks.Validate()
	assert.NotNil(err)
}

func TestKSTestConfValidateShouldFailBadNormalCount(t *testing.T) {
	assert := assert.New(t)

	ks := kstest.KSTestConf{NormalCount: 0}

	err := ks.Validate()
	assert.NotNil(err)
}

func TestKSTestConfValidateShouldFailBadMaxFingerprints(t *testing.T) {
	assert := assert.New(t)

	ks := kstest.KSTestConf{MaxFingerprints: 0}

	err := ks.Validate()
	assert.NotNil(err)
}

func TestKSTestConfValidateShouldPass(t *testing.T) {
	assert := assert.New(t)

	ks := kstest.KSTestConf{
		Strictness:      1,
		NormalCount:     3,
		MaxFingerprints: 20,
	}

	err := ks.Validate()
	assert.Nil(err)
}

func TestKSTestConfShouldParse(t *testing.T) {
	assert := assert.New(t)

	var data string = `---
strictness: 4
normal_count: 4
max_fingerprints: 10
`

	ks := kstest.KSTestConf{}

	err := yaml.Unmarshal([]byte(data), &ks)

	assert.Nil(err)

	assert.Equal(4, ks.Strictness)
	assert.Equal(4, ks.NormalCount)
	assert.Equal(10, ks.MaxFingerprints)
}
