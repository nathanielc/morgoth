package engine_test

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nvcook42/morgoth/engine"
	"testing"
)

func TestEngineShouldHaveRegistery(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(engine.Registery)
}
