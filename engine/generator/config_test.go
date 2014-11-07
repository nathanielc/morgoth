package generator_test

import (
	"github.com/nvcook42/morgoth/engine/generator"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestGeneratorConfShouldDefaultEmpty(t *testing.T) {
	assert := assert.New(t)

	conf := generator.GeneratorConf{}
	conf.Default()

	assert.Equal(conf.Host, "localhost")
	assert.Equal(conf.Port, 8083)

}

func TestGeneratorConfShouldDefaultNonEmpty(t *testing.T) {
	assert := assert.New(t)

	conf := generator.GeneratorConf{
		Port:     65536,
		Database: "morgoth",
	}
	conf.Default()

	assert.Equal(conf.Host, "localhost")
	assert.Equal(conf.Port, 8083)
	assert.Equal(conf.Database, "morgoth")

}

func TestGeneratorConfDefaultShouldIgnoreValidFields(t *testing.T) {
	assert := assert.New(t)

	conf := generator.GeneratorConf{
		Host:     "influx",
		Port:     42,
		Database: "morgoth",
	}
	conf.Default()

	assert.Equal(conf.Host, "influx")
	assert.Equal(conf.Port, 42)
	assert.Equal(conf.Database, "morgoth")

}

func TestGeneratorConfShouldParse(t *testing.T) {
	assert := assert.New(t)

	ic := generator.GeneratorConf{}

	var data string = `---
host: influx1.example.com
port: 4242
user: morgoth
password: mysecret
database: morgothdb
`

	err := yaml.Unmarshal([]byte(data), &ic)

	assert.Nil(err)

	assert.Equal("influx1.example.com", ic.Host)
	assert.Equal(4242, ic.Port)
	assert.Equal("morgoth", ic.User)
	assert.Equal("mysecret", ic.Password)
	assert.Equal("morgothdb", ic.Database)

}
