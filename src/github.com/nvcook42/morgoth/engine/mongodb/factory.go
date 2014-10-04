package mongodb

import (
	"github.com/nvcook42/morgoth/engine"
	"github.com/nvcook42/morgoth/config/types"
)

type MongoDBFactory struct {
}

func (self *MongoDBFactory) NewConf() types.Configuration {
	return new(MongoDBConf)
}

func (self *MongoDBFactory) GetInstance(config types.Configuration) (interface{}, error) {
	return new(MongoDBEngine), nil
}

func init() {
	factory := new(MongoDBFactory)
	engine.Registery.RegisterFactory("mongodb", factory)
}
