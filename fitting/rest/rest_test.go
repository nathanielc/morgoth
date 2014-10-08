package rest_test

import (
	"github.com/nvcook42/morgoth/fitting"
	_ "github.com/nvcook42/morgoth/fitting/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRestShouldBeRegistered(t *testing.T) {
	assert := assert.New(t)

	_, err := fitting.Registery.GetFactory("rest")

	assert.Nil(err)
}
