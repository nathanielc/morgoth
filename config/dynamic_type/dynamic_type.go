// This package provides helper methods for parsing YAML configuration
// based on a dynamic type of configuration
package dynamic_type

import (
	"errors"
	"fmt"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/nathanielc/morgoth/config/types"
	"github.com/nathanielc/morgoth/registery"
)

type configContainer interface {
	GetConf() types.Configuration
}

// Base object for unmarshaling dyanmic yaml
// into a Configuration object.
type DynamicConfiguration struct {
	Type      string
	Conf      types.Configuration
	registery *registery.Registery
}

func (self *DynamicConfiguration) Default() {
	if self.Conf != nil {
		self.Conf.Default()
	}
}

func (self DynamicConfiguration) Validate() error {
	if self.Conf == nil {
		return errors.New("Conf is nil, cannot validate")
	}
	return self.Conf.Validate()
}

func (self *DynamicConfiguration) GetConf() types.Configuration {
	return self.Conf
}

func PerformFromYAML(yamlConf string, confStruct configContainer) error {
	yaml.Unmarshal([]byte(yamlConf), confStruct)
	conf := confStruct.GetConf()
	conf.Default()
	err := conf.Validate()
	if err != nil {
		return err
	}
	return nil
}

// Return instance via the Registery from the parsed config
func (self *DynamicConfiguration) PerformGetInstance(registery *registery.Registery) (interface{}, error) {
	factory, err := registery.GetFactory(self.Type)
	if err != nil {
		return nil, err
	}

	instance, err := factory.GetInstance(self.Conf)
	if err != nil {
		return nil, err
	}
	if instance == nil {
		return nil, errors.New("Instance is nil")
	}

	return instance, nil
}

// Performs the unmarshaling into self given a registery and
// the unmarshal function.
func (self *DynamicConfiguration) PerformUnmarshalYAML(registery *registery.Registery, unmarshal func(interface{}) error) error {
	t, config, err := UnmarshalDynamicType(registery, unmarshal)
	if err != nil {
		return err
	}
	self.Type = t
	self.Conf = config
	return nil
}

// Unmarshal YAML configuration into a dynamically named configuration struct
// using a Registery.
//
// Consider the following yaml
//
//    mysql:
//      host: localhost
//      port: 3307
//
// The 'mysql' key will be read as the 'type'. Then using a Registery,
// a factory will be used to create a new empty configuration of type 'mysql'.
// The remaining yaml will be unmarshaled into the 'mysql' configuration struct.
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

///////////////////////////
// Internals
///////////////////////////

var (
	conf types.Configuration
)

type confUnmarshaler struct {
}

func (self *confUnmarshaler) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(conf)
}
