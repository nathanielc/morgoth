package mocks

import "github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/mock"

type Configuration struct {
	mock.Mock
}

func (m *Configuration) Default() {
	m.Called()
}

func (m *Configuration) Validate() error {
	m.Called()
	return nil
}
