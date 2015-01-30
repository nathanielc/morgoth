package detector_test

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nvcook42/morgoth/detector"
	"testing"
)

func TestDetectorShouldHaveRegistery(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(detector.Registery)

}
