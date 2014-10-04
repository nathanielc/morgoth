package dynamic_type

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/nvcook42/morgoth/registery"
	"gopkg.in/yaml.v2"
)

type testStruct struct {
	assert *assert.Assertions
	registery *registery.Registery
}

type testConfig struct {
	A int
	B int
	C int
}

func (self *testConfig) Default() {
}

func (self testConfig) Validate() error {
	return nil
}


type testFactory struct {
}

func (self *testFactory) NewConf() registery.Configuration {
	return new(testConfig)
}

func (self *testFactory) GetInstance(config registery.Configuration) (interface{}, error) {
	return nil, nil
}

func (self *testStruct) UnmarshalYAML(unmarshal func(interface{}) error) error {
	engineType, config, err := UnmarshalDynamicType(self.registery, unmarshal)
	self.assert.Nil(err)
	self.assert.Equal("jim", engineType)
	if !self.assert.NotNil(config) {
		self.assert.Fail("Config was nil")
	}
	self.assert.Equal(testConfig{1, 2, 4}, *config.(*testConfig))
	return nil
}

func TestDynamicType(t *testing.T) {
	assert := assert.New(t)
	registery := registery.New()
	tf := testFactory{}
	registery.RegisterFactory("jim", &tf)

	ts := testStruct{assert, registery}
	var data = `
jim:
  a: 1
  b: 2
  c: 4
`
	err := yaml.Unmarshal([]byte(data), &ts)
	assert.Nil(err)
	

	

}
