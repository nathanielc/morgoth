package registery_test

import (
	mocks "github.com/nvcook42/morgoth/mocks/registery"
	"github.com/nvcook42/morgoth/registery"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegistryShouldRegister(t *testing.T) {
	assert := assert.New(t)

	r := registery.New()

	mockFactory := new(mocks.Factory)

	err := r.RegisterFactory("mocktest", mockFactory)
	assert.Nil(err)

	factory, err := r.GetFactory("mocktest")
	assert.Nil(err)

	assert.Equal(mockFactory, factory)

	mockFactory.Mock.AssertExpectations(t)

}
