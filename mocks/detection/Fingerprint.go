package mocks

import "github.com/nathanielc/morgoth/detection"
import "github.com/stretchr/testify/mock"

type Fingerprint struct {
	mock.Mock
}

func (m *Fingerprint) IsMatch(other detection.Fingerprint) bool {
	ret := m.Called(other)

	r0 := ret.Get(0).(bool)

	return r0
}
