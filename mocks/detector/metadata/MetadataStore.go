package mocks

import "github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/mock"

import metric "github.com/nvcook42/morgoth/metric/types"

type MetadataStore struct {
	mock.Mock
}

func (m *MetadataStore) StoreDoc(_a0 metric.MetricID, _a1 []byte) {
	m.Called(_a0, _a1)
}
func (m *MetadataStore) GetDoc(_a0 metric.MetricID) []byte {
	ret := m.Called(_a0)

	r0 := ret.Get(0).([]byte)

	return r0
}
func (m *MetadataStore) Close() {
	m.Called()
}
