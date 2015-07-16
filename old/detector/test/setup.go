package test

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/mock"
	config "github.com/nathanielc/morgoth/config/types"
	"github.com/nathanielc/morgoth/detector"
	"github.com/nathanielc/morgoth/engine/generator"
	metric "github.com/nathanielc/morgoth/metric/types"
	app "github.com/nathanielc/morgoth/mocks/app"
	metadata "github.com/nathanielc/morgoth/mocks/detector/metadata"
	"github.com/nathanielc/morgoth/registery"
	"github.com/nathanielc/morgoth/schedule"
)

func SetupGeneratedDetectorTest(
	metric metric.MetricID,
	rotation schedule.Rotation,
	conf config.Configuration,
	factory registery.Factory,
	functions map[metric.MetricID]generator.Ft,
) (detector.Detector, error) {
	conf.Default()
	err := conf.Validate()
	if err != nil {
		return nil, err
	}

	d, err := factory.GetInstance(conf)
	if err != nil {
		return nil, err
	}

	detector := d.(detector.Detector)

	engine := generator.New(functions)
	engine.Initialize()

	app := new(app.App)
	meta := new(metadata.MetadataStore)

	meta.On("GetDoc", mock.Anything).Return([]byte{})
	meta.On("StoreDoc", mock.Anything, mock.Anything).Return()

	app.On("GetWriter").Return(engine).Once()
	app.On("GetReader").Return(engine).Once()
	app.On("GetMetadataStore", mock.Anything).Return(meta, nil).Once()

	err = detector.Initialize(app, rotation)
	if err != nil {
		return nil, err
	}

	return detector, nil
}
