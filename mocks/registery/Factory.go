package mocks

import "github.com/stretchr/testify/mock"

import "github.com/nvcook42/morgoth/config/types"

type Factory struct {
	mock.Mock
}

func (m *Factory) NewConf() types.Configuration {
	m.Called()
	return nil
}
func (m *Factory) GetInstance(_a0 types.Configuration) (interface{}, error) {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(interface{})
	r1 := ret.Error(1)

	return r0, r1
}
