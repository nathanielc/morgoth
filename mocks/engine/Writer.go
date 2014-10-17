package mocks

import "github.com/stretchr/testify/mock"

import metric "github.com/nvcook42/morgoth/metric/types"
import "time"

type Writer struct {
 mock.Mock
}

func (m *Writer) Insert(datetime time.Time, metric metric.MetricID, value float64) {
 m.Called(datetime, metric, value)
}
