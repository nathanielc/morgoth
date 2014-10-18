package fileio

import (
	"github.com/nvcook42/morgoth/config/types"
	"github.com/nvcook42/morgoth/engine"
	"errors"
	"fmt"
)

type FileIOFactory struct {
}

func (self *FileIOFactory) NewConf() types.Configuration {
	return new(FileIOConf)
}

func (self *FileIOFactory) GetInstance(config types.Configuration) (interface{}, error) {
	fileIOConf, ok := config.(*FileIOConf)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Configuration is not FileIOConf%v", config))
	}
	engine := &FileIOEngine{
		Dir: fileIOConf.Dir,
	}
	return engine, nil
}

func init() {
	factory := new(FileIOFactory)
	engine.Registery.RegisterFactory("fileio", factory)
}
