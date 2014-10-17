package mocks

import "github.com/nvcook42/morgoth/engine"
import "github.com/stretchr/testify/mock"

type Engine struct {
 mock.Mock
}

func (m *Engine) Initialize() error {
 ret := m.Called()

 r0 := ret.Error(0)

 return r0
}
func (m *Engine) GetReader() engine.Reader {
 ret := m.Called()

 r0 := ret.Get(0).(engine.Reader)

 return r0
}
func (m *Engine) GetWriter() engine.Writer {
 ret := m.Called()

 r0 := ret.Get(0).(engine.Writer)

 return r0
}
