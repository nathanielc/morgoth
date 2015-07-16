package mocks

import "github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/mock"

import "github.com/nathanielc/morgoth/engine"
import "github.com/nathanielc/morgoth/schedule"
import "github.com/nathanielc/morgoth/detector/metadata"

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
func (m *App) GetMetadataStore(detectorID string) (metadata.MetadataStore, error) {
	ret := m.Called(detectorID)

	r0 := ret.Get(0).(metadata.MetadataStore)
	r1 := ret.Error(1)

	return r0, r1
}
