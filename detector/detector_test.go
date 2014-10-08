package detector_test

import (
	"github.com/nvcook42/morgoth/detector"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDetectorShouldHaveRegistery(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(detector.Registery)

}
