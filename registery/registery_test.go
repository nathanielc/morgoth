package registery_test

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	mocks "github.com/nathanielc/morgoth/mocks/registery"
	"github.com/nathanielc/morgoth/registery"
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
