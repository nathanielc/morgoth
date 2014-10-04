package dynamic_type

import (
	"errors"
	"github.com/nvcook42/morgoth/registery"
)

var (
	conf registery.Configuration
)


type confUnmarshaler struct {
}

func (self *confUnmarshaler) UnmarshalYAML(unmarshal func(interface{}) error ) error {
	return unmarshal(conf)
}


func UnmarshalDynamicType(
	reg *registery.Registery,
	unmarshal func(interface{}) error,
) (string, registery.Configuration, error) {

	typeData := make(map[string]interface{})
	err := unmarshal(&typeData)
	if err != nil {
		return "", nil, err
	}

	if len(typeData) != 1 {
		return "", nil, errors.New("Only one key can be specified")
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


