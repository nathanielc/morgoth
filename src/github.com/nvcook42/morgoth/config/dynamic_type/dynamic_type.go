package dynamic_type

import (
	"errors"
	"fmt"
	"github.com/nvcook42/morgoth/registery"
	"github.com/nvcook42/morgoth/config/types"
)

var (
	conf types.Configuration
)


type confUnmarshaler struct {
}

func (self *confUnmarshaler) UnmarshalYAML(unmarshal func(interface{}) error ) error {
	return unmarshal(conf)
}


func UnmarshalDynamicType(
	reg *registery.Registery,
	unmarshal func(interface{}) error,
) (string, types.Configuration, error) {

	typeData := make(map[string]interface{})
	err := unmarshal(&typeData)
	if err != nil {
		return "", nil, err
	}

	if len(typeData) != 1 {
		return "", nil, errors.New(fmt.Sprintf("Exactly one key must be specified. Found: %v", typeData))
	}
	var typeName string
	for key := range typeData {
		typeName = key
	}

	factory, err := reg.GetFactory(typeName)
	if err != nil {
		return "", nil, err
	}

	confData := make(map[string]confUnmarshaler)
	conf = factory.NewConf()

	err = unmarshal(&confData)
	if err != nil {
		return "", nil, err
	}

	return typeName, conf, nil
}


