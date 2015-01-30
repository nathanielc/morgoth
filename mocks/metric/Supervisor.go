package mocks

import "github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/mock"

import "github.com/nvcook42/morgoth/metric/types"

import "github.com/nvcook42/morgoth/schedule"
import "time"

type Supervisor struct {
	mock.Mock
}

func (m *Supervisor) GetPattern() types.Pattern {
	ret := m.Called()

	r0 := ret.Get(0).(types.Pattern)

	return r0
}
func (m *Supervisor) AddMetric(_a0 types.MetricID) {
	m.Called(_a0)
}
func (m *Supervisor) Detect(rotation schedule.Rotation, start time.Time, stop time.Time) {
	m.Called(rotation, start, stop)
}
