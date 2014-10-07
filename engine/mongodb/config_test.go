package mongodb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMongoDBConfShouldDefaultEmpty(t *testing.T) {
	assert := assert.New(t)

	conf := MongoDBConf{}
	conf.Default()

	assert.Equal(conf.Host, "localhost")
	assert.Equal(conf.Port, 27017)

}

func TestMongoDBConfShouldDefaultNonEmpty(t *testing.T) {
	assert := assert.New(t)

	conf := MongoDBConf{
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

	conf := MongoDBConf{
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
