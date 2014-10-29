package mocks

import "github.com/stretchr/testify/mock"

import "github.com/nvcook42/morgoth/metric/types"
import app "github.com/nvcook42/morgoth/app/types"

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
func (m *Supervisor) Start(_a0 app.App) {
	m.Called(_a0)
}
