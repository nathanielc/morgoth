package mongodb

import (
	"github.com/nvcook42/morgoth/engine"
	"github.com/nvcook42/morgoth/registery"
)

type MongoDBFactory struct {
}

func (self *MongoDBFactory) NewConf() registery.Configuration {
	return new(MongoDBConf)
}

func (self *MongoDBFactory) GetInstance(config registery.Configuration) (interface{}, error) {
	return new(MongoDBEngine), nil
}

func init() {
	factory := new(MongoDBFactory)
	engine.Registery.RegisterFactory("mongodb", factory)
}
