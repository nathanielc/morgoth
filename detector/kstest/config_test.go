package kstest_test

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/nathanielc/morgoth/detector/kstest"
	"testing"
)

func TestKSTestConfShouldDefault(t *testing.T) {
	assert := assert.New(t)

	ks := kstest.KSTestConf{}

	ks.Default()

	assert.Equal(1, ks.Confidence)
	assert.Equal(3, ks.NormalCount)
	assert.Equal(20, ks.MaxFingerprints)

}

func TestKSTestConfValidateShouldFailBadConfidence(t *testing.T) {
	assert := assert.New(t)

	ks := kstest.KSTestConf{Confidence: 6}

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
		Confidence:      1,
		NormalCount:     3,
		MaxFingerprints: 20,
	}

	err := ks.Validate()
	assert.Nil(err)
}

func TestKSTestConfShouldParse(t *testing.T) {
	assert := assert.New(t)

	var data string = `---
confidence: 4
normal_count: 4
max_fingerprints: 10
`

	ks := kstest.KSTestConf{}

	err := yaml.Unmarshal([]byte(data), &ks)

	assert.Nil(err)

	assert.Equal(4, ks.Confidence)
	assert.Equal(4, ks.NormalCount)
	assert.Equal(10, ks.MaxFingerprints)
}
