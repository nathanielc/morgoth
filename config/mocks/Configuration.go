package mocks

import "github.com/stretchr/testify/mock"

type Configuration struct {
	mock.Mock
}

func (m *Configuration) Validate() error {
	ret := m.Called()

	r0 := ret.Error(0)

	return r0
}

func (m *Configuration) Default() {
	m.Called()
}
