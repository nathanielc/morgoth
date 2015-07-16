package detector

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"

	"testing"
)

type mockFingerprint struct {
	id int
}


	eq := self == other
	return eq
}

type mockFingerprinter struct{}


	return mockFingerprint{int(window[0])}
}

func TestDetectorShouldHonorNormalCount(t *testing.T) {
	assert := assert.New(t)

	normalCount := 3
	consensus := 0.5
	minSupport := 0.1
	errorTolerance := 0.01


	fingerprinter := &mockFingerprinter{}
	d.AddFingerprinter(fingerprinter)

	assert.True(d.IsAnomalous([]float64{1}))
	assert.True(d.IsAnomalous([]float64{1}))
	assert.False(d.IsAnomalous([]float64{1}))
	assert.False(d.IsAnomalous([]float64{1}))
	assert.True(d.IsAnomalous([]float64{2}))

	normalCount = 2


	d.AddFingerprinter(fingerprinter)

	assert.True(d.IsAnomalous([]float64{1}))
	assert.False(d.IsAnomalous([]float64{1}))
	assert.False(d.IsAnomalous([]float64{1}))
}
