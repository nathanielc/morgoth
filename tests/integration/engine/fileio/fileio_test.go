package fileio_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/nvcook42/morgoth/engine/fileio"
	"github.com/nvcook42/morgoth/tests/integration/engine"
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


