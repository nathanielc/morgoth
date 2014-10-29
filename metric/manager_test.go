package metric_test

import (
	"github.com/nvcook42/morgoth/metric"
	"github.com/nvcook42/morgoth/metric/types"
	app "github.com/nvcook42/morgoth/mocks/app"
	mock "github.com/nvcook42/morgoth/mocks/metric"
	"testing"
)

func TestManagerShouldHandleNewMetric(t *testing.T) {

	app := new(app.App)

	var metricName types.MetricID = "metricName"
	ms0 := new(mock.Supervisor)
	ms0.On("GetPattern").Return(types.Pattern(".*")).Once()
	ms0.On("AddMetric", metricName).Return().Once()
	ms0.On("Start", app).Return().Once()

	ms1 := new(mock.Supervisor)
	ms1.On("GetPattern").Return(types.Pattern("\\d+")).Once()

	ms2 := new(mock.Supervisor)
	ms2.On("GetPattern").Return(types.Pattern(".*")).Once()


	mockSupervisors := []*mock.Supervisor{
		ms0,
		ms1,
		ms2,
	}
	supervisors := make([]metric.Supervisor, len(mockSupervisors))
	for i := range mockSupervisors {
		supervisors[i] = mockSupervisors[i]
	}


	m := metric.NewManager(app, supervisors)
	m.NewMetric(metricName)
	m.NewMetric(metricName)

	for i := range mockSupervisors {
		mockSupervisors[i].Mock.AssertExpectations(t)
	}

	app.AssertExpectations(t)

}
