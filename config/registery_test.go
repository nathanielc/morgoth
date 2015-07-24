package config_test

//import (
//	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
//	"github.com/nathanielc/morgoth/config"
//	"testing"
//)
//
//func TestRegistryShouldRegister(t *testing.T) {
//	assert := assert.New(t)
//
//	r := config.NewRegistry()
//
//	mockFactory := new(config.Factory)
//
//	err := r.RegisterFactory("mocktest", mockFactory)
//	assert.Nil(err)
//
//	factory, err := r.GetFactory("mocktest")
//	assert.Nil(err)
//
//	assert.Equal(mockFactory, factory)
//
//	mockFactory.Mock.AssertExpectations(t)
//
//}
