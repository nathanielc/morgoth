package detectors_test

import (
	"testing"
	"github.com/golang/glog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	app "github.com/nvcook42/morgoth/mocks/app"
	metadata "github.com/nvcook42/morgoth/mocks/detector/metadata"
	"github.com/nvcook42/morgoth/schedule"
	"github.com/nvcook42/morgoth/engine/generator"
	"github.com/nvcook42/morgoth/detector/mgof"
	"time"
)

var rotation = schedule.Rotation{
	Period: time.Minute,
	Resolution: time.Second,
}

func TestMGOF(t *testing.T) {
	defer glog.Flush()
	assert := assert.New(t)

	factory := mgof.MGOFFactory{}

	mgofConf := &mgof.MGOFConf{
		Min: 0,
		Max: 100,
	}
	mgofConf.Default()
	
	detector, err := factory.GetInstance(mgofConf)
	if !assert.Nil(err) {
		assert.Fail("Failed to create mgof detector")
	}
	mgofDetector, ok := detector.(*mgof.MGOF)
	if !assert.True(ok) {
		assert.Fail("Detector not a MGOF detector", mgofDetector)
	}

	engine := new(generator.GeneratorEngine)
	engine.Initialize()
	
	app := new(app.App)
	meta := new(metadata.MetadataStore)


	meta.On("GetDoc", "m1").Return([]byte{})
	meta.On("StoreDoc", "m1", mock.Anything).Return()

	app.On("GetWriter").Return(engine).Once()
	app.On("GetReader").Return(engine).Once()
	app.On("GetMetadataStore", mock.Anything).Return(meta, nil).Once()

	err = mgofDetector.Initialize(app, rotation)
	assert.Nil(err)

	mgof.TraceLevel = -1

	start := time.Time{}.Add(time.Second) //Start one second past zero
	stop := start.Add(rotation.Period)
	anomalous := mgofDetector.Detect("m1", start, stop)
	assert.True(anomalous)

}
