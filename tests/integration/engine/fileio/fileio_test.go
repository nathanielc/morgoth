package fileio_test

import (
	"github.com/nvcook42/morgoth/engine/fileio"
	"github.com/nvcook42/morgoth/tests/integration/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAll(t *testing.T) {
	assert := assert.New(t)
	fileioEngine := fileio.FileIOEngine{
		Dir: "/tmp/morgoth/fileiodb/",
	}
	err := fileioEngine.Initialize()
	assert.Nil(err)
	suite := engine.EngineTestSuite{&fileioEngine}
	suite.TestAll(t)
}
