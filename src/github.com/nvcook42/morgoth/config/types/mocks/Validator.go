package mocks

import "github.com/stretchr/testify/mock"

type Validator struct {
 mock.Mock
}

func (m *Validator) Validate() error {
 ret := m.Called()

 r0 := ret.Error(0)

 return r0
}
