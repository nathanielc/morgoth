package config

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func LoadFromFile(path string) (*Config, error) {
	config_file, err := os.Open(path)
	if err != nil {
		return nil, err

	}
	data, err := ioutil.ReadAll(config_file)
	if err != nil {
		return nil, err
	}
	return Load(data)
}

func Load(data []byte) (*Config, error) {
	config := new(Config)

	err := yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	config.Default()
	err = config.Validate()
	if err != nil {
		return nil, err
	}

	return config, nil
}
