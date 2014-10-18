package engine_test

import (
	"github.com/nvcook42/morgoth/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEngineShouldHaveRegistery(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(engine.Registery)
}

