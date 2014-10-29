package mocks

import "github.com/stretchr/testify/mock"

import "github.com/nvcook42/morgoth/engine"
import "github.com/nvcook42/morgoth/schedule"

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
func (m *App) GetSchedule() schedule.Schedule {
	ret := m.Called()

	r0 := ret.Get(0).(schedule.Schedule)

	return r0
}
