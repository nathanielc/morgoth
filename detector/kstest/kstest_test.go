package kstest_test

import (
	"github.com/nvcook42/morgoth/detector"
	_ "github.com/nvcook42/morgoth/detector/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKSTestShouldBeRegistered(t *testing.T) {
	assert := assert.New(t)

	_, err := detector.Registery.GetFactory("kstest")
	assert.Nil(err)
}
