package detection_test

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nvcook42/morgoth/detection"
	"testing"
)

type mockFingerprint struct {
	id int
}

func (self mockFingerprint) IsMatch(other detection.Fingerprint) bool {
	eq := self == other
	return eq
}

type mockFingerprinter struct{}

func (self *mockFingerprinter) Fingerprint(window []float64) detection.Fingerprint {
	return mockFingerprint{int(window[0])}
}

func TestDetectorShouldHonorNormalCount(t *testing.T) {
	assert := assert.New(t)

	normalCount := 3
	consensus := 0.5
	minSupport := 0.1
	errorTolerance := 0.01
	d := detection.New(normalCount, consensus, minSupport, errorTolerance)

	fingerprinter := &mockFingerprinter{}
	d.AddFingerprinter(fingerprinter)

	assert.True(d.IsAnomalous([]float64{1}))
	assert.True(d.IsAnomalous([]float64{1}))
	assert.False(d.IsAnomalous([]float64{1}))
	assert.False(d.IsAnomalous([]float64{1}))
	assert.True(d.IsAnomalous([]float64{2}))

	normalCount = 2
	d = detection.New(normalCount, consensus, minSupport, errorTolerance)

	d.AddFingerprinter(fingerprinter)

	assert.True(d.IsAnomalous([]float64{1}))
	assert.False(d.IsAnomalous([]float64{1}))
	assert.False(d.IsAnomalous([]float64{1}))
}
