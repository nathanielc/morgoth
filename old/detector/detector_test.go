package detector_test

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nathanielc/morgoth/detector"
	"testing"
)

func TestDetectorShouldHaveRegistery(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(detector.Registery)

}
