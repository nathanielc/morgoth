package config

import (
	"errors"
	"github.com/nvcook42/morgoth/data/influxdb"
	"github.com/nvcook42/morgoth/data/mongodb"
)

// Base config struct for the entire morgoth config
type Config struct {
	DataEngine DataEngine `yaml:"data_engine"`
	Metrics    Metrics    `yaml:"metrics"`
	Fittings   Fittings   `yaml:"fittings"`
}

type EngineType string
func (self EngineType) Validate() error {
	if self == InfluxDB || self == MongoDB {
		return nil
	}
	return errors.New("Invalid EngineType")
}

const (
	InfluxDB EngineType = "influxdb"
	MongoDB  EngineType = "mongodb"
)


// The Data Engine subsection of the config
type DataEngine struct {
	Type     EngineType            `yaml:"type"`
	InfluxDB influxdb.InfluxDBConf `yaml:"influxdb,omitempty"`
	MongoDB  mongodb.MongoDBConf   `yaml:"mongodb,omitempty"`
}

// The Metrics subsection of the config
type Metrics struct {
	Metrics []MetricConf `yaml:"metrics"`
}

// Represents a single metric conf section
type MetricConf struct {
	Pattern  string   `yaml:"pattern"`
	Schedule Schedule `yaml:"schedule"`
}

type Schedule struct {
	Duration uint
	Period   uint
	Delay    uint
}

// The Fittings subsection of the config
type Fittings struct {
}
