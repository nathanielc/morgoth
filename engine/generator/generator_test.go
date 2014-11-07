package generator_test

import (
	"github.com/nvcook42/morgoth/engine"
	_ "github.com/nvcook42/morgoth/engine/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneratorShouldBeRegistered(t *testing.T) {
	assert := assert.New(t)

	_, err := engine.Registery.GetFactory("generator")
	assert.Nil(err)
}
