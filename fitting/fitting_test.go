package fitting_test

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nvcook42/morgoth/fitting"
	"testing"
)

func TestFittingShouldHaveRegistery(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(fitting.Registery)
}
