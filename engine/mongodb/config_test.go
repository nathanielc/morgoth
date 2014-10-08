package mongodb_test

import (
	"github.com/nvcook42/morgoth/engine/mongodb"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestMongoDBConfShouldDefaultEmpty(t *testing.T) {
	assert := assert.New(t)

	conf := mongodb.MongoDBConf{}
	conf.Default()

	assert.Equal(conf.Host, "localhost")
	assert.Equal(conf.Port, 27017)

}

func TestMongoDBConfShouldDefaultNonEmpty(t *testing.T) {
	assert := assert.New(t)

	conf := mongodb.MongoDBConf{
		Port:     65536,
		Database: "morgoth",
	}
	conf.Default()

	assert.Equal(conf.Host, "localhost")
	assert.Equal(conf.Port, 27017)
	assert.Equal(conf.Database, "morgoth")

}

func TestMongoDBConfDefaultShouldIgnoreValidFields(t *testing.T) {
	assert := assert.New(t)

	conf := mongodb.MongoDBConf{
		Host:      "mongo",
		Port:      42,
		Database:  "morgoth",
		IsSharded: true,
	}
	conf.Default()

	assert.Equal(conf.Host, "mongo")
	assert.Equal(conf.Port, 42)
	assert.Equal(conf.Database, "morgoth")
	assert.Equal(conf.IsSharded, true)

}

func TestMongoDBConfShouldParse(t *testing.T) {
	assert := assert.New(t)

	mc := mongodb.MongoDBConf{}

	var data string = `---
host: mongo1.example.com
port: 4242
database: morgothdb
is_sharded: true
`

	err := yaml.Unmarshal([]byte(data), &mc)

	assert.Nil(err)

	assert.Equal("mongo1.example.com", mc.Host)
	assert.Equal(4242, mc.Port)
	assert.Equal("morgothdb", mc.Database)
	assert.True(mc.IsSharded)

}
