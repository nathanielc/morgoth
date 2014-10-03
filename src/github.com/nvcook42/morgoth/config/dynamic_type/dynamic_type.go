package dynamic_type

import (
	"github.com/nvcook42/morgoth/registery"
	log "github.com/cihub/seelog"
)

type typeshell struct {
	Type string `yaml:"type"`
}


func UnmarshalDynamicType(
	key string,
	reg *registery.Registery,
	unmarshal func(interface{}) error,
) (string, registery.Configuration, error) {
	defer log.Flush()

	ts := typeshell{}
	err := unmarshal(&ts)
	if err != nil {
		return "", nil, err
	}

	log.Debugf("dynamic_type Type: %s", ts.Type)

	factory, err := reg.GetFactory(ts.Type)
	if err != nil {
		return "", nil, err
	}

	data := make(map[string]registery.Configuration)
	data[key] = factory.NewConf()

	err = unmarshal(&data)
	if err != nil {
		return "", nil, err
	}

	log.Debugf("key: %s type: %s, config: %v", key, ts.Type, data)
	return ts.Type, data[key].(registery.Configuration), nil
}


