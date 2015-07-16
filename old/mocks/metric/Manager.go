package mocks

import "github.com/nathanielc/morgoth/metric/types"
import "github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/mock"

import "github.com/nathanielc/morgoth/schedule"
import "time"

type Manager struct {
	mock.Mock
}

func (m *Manager) NewMetric(_a0 types.MetricID) {
	m.Called(_a0)
}
func (m *Manager) Detect(rotation schedule.Rotation, start time.Time, stop time.Time) {
	m.Called(rotation, start, stop)
}
