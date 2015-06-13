package engine_test

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nathanielc/morgoth/engine"
	"testing"
)

func TestEngineShouldHaveRegistery(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(engine.Registery)
}
