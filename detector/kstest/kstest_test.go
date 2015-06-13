package kstest_test

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nathanielc/morgoth/detector"
	_ "github.com/nathanielc/morgoth/detector/list"
	"testing"
)

func TestKSTestShouldBeRegistered(t *testing.T) {
	assert := assert.New(t)

	_, err := detector.Registery.GetFactory("kstest")
	assert.Nil(err)
}
