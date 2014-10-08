package mongodb_test

import (
	"github.com/nvcook42/morgoth/engine"
	_ "github.com/nvcook42/morgoth/engine/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMongoDBShouldBeRegistered(t *testing.T) {
	assert := assert.New(t)

	_, err := engine.Registery.GetFactory("mongodb")
	assert.Nil(err)
}
