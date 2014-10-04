package mongodb

import (
	"github.com/nvcook42/morgoth/engine"
)

type MongoDBEngine struct {
}

func (self *MongoDBEngine) GetReader() engine.Reader {
	return nil
}

func (self *MongoDBEngine) GetWriter() engine.Writer {
	return nil
}
