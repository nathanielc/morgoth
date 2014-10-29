package mocks

import "github.com/nvcook42/morgoth/metric/types"
import "github.com/stretchr/testify/mock"

type Manager struct {
	mock.Mock
}

func (m *Manager) NewMetric(_a0 types.MetricID) {
	m.Called(_a0)
}
