package mocks

import "github.com/stretchr/testify/mock"

import "github.com/nvcook42/morgoth/engine"
import metric "github.com/nvcook42/morgoth/metric/types"

type App struct {
 mock.Mock
}

func (m *App) Run() error {
 ret := m.Called()

 r0 := ret.Error(0)

 return r0
}
func (m *App) GetWriter() engine.Writer {
 ret := m.Called()

 r0 := ret.Get(0).(engine.Writer)

 return r0
}
func (m *App) GetReader() engine.Reader {
 ret := m.Called()

 r0 := ret.Get(0).(engine.Reader)

 return r0
}
func (m *App) GetManager() metric.Manager {
 ret := m.Called()

 r0 := ret.Get(0).(metric.Manager)

 return r0
}
