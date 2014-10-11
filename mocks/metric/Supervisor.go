package mocks

import "github.com/nvcook42/morgoth/metric"
import "github.com/stretchr/testify/mock"

import "github.com/nvcook42/morgoth/metric/types"

type Supervisor struct {
	mock.Mock
}

func (m *Supervisor) GetPattern() metric.Pattern {
	ret := m.Called()

	r0 := ret.Get(0).(metric.Pattern)

	return r0
}
func (m *Supervisor) AddMetric(_a0 types.MetricID) {
	m.Called(_a0)
}
func (m *Supervisor) Start() {
	m.Called()
}
