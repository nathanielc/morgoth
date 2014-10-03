package dynamic_type

import (
	log "github.com/cihub/seelog"
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
	engineType, config, err := UnmarshalDynamicType("asdf", self.registery, unmarshal)
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
type: jim
asdf:
  a: 1
  b: 2
  c: 4
`
	err := yaml.Unmarshal([]byte(data), &ts)
	assert.Nil(err)
	

	

}

func TestYaml(t *testing.T) {
	assert := assert.New(t)

var data = `
asdf:
  a: 1
  b: 2
  c: 4
qwerty:
  a: 5
  b: 8
  c: 6
`
	type s struct {
		data map[string]testConfig
	}
	ts := new(s)
	ts.data = make(map[string]testConfig)
	err := yaml.Unmarshal([]byte(data), &ts.data)
	assert.Nil(err)

	log.Debugf("TS: %v", *ts)
	

}
