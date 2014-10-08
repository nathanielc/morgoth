package fitting_test

import (
	"github.com/nvcook42/morgoth/fitting"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFittingShouldHaveRegistery(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(fitting.Registery)
}
