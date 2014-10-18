package fileio_test


import (
	"github.com/nvcook42/morgoth/engine/fileio"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestFileIOConfShouldDefaultEmpty(t *testing.T) {
	assert := assert.New(t)

	conf := fileio.FileIOConf{}
	conf.Default()

	assert.Equal("./fileiodb/", conf.Dir)
}

func TestFileIOConfShouldDefaultNonEmpty(t *testing.T) {
	assert := assert.New(t)

	conf := fileio.FileIOConf{
		Dir: "/tmp",
	}
	conf.Default()

	assert.Equal("/tmp", conf.Dir)

}

func TestFileIOConfShouldParse(t *testing.T) {
	assert := assert.New(t)

	mc := fileio.FileIOConf{}

	var data string = `---
dir: /tmp/morgoth/fileiodb/
`

	err := yaml.Unmarshal([]byte(data), &mc)

	assert.Nil(err)

	assert.Equal("/tmp/morgoth/fileiodb/", mc.Dir)

}
