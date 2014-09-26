package config

import (
	"errors"
	"github.com/nvcook42/morgoth/data/influxdb"
	"github.com/nvcook42/morgoth/data/mongodb"
	"github.com/nvcook42/morgoth/validator"
)

// Supported EngineTypes
type EngineType string

const (
	InfluxDB EngineType = "influxdb"
	MongoDB  EngineType = "mongodb"
)

func (self EngineType) Validate() error {
	switch self {
	case
		InfluxDB,
		MongoDB:
		return nil
	}
	return errors.New("Invalid EngineType")
}

// The Data Engine subsection of the config
type DataEngine struct {
	Type     EngineType             `yaml:"type"`
	InfluxDB *influxdb.InfluxDBConf `yaml:"influxdb,omitempty"`
	MongoDB  *mongodb.MongoDBConf   `yaml:"mongodb,omitempty"`
}

func (self *DataEngine) Default() {
	switch self.Type {
	case InfluxDB:
		self.InfluxDB.Default()
	case MongoDB:
		self.MongoDB.Default()
	}
}

func (self DataEngine) Validate() error {
	return validator.ValidateOne(self)
}
