package fitting_test

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nathanielc/morgoth/fitting"
	"testing"
)

func TestFittingShouldHaveRegistery(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(fitting.Registery)
}
