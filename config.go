package morgoth

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/nathanielc/morgoth/config"
	"io/ioutil"
	"os"
	"time"
)

// Base config struct for the entire morgoth config
type Config struct {
	APIPort    uint            `yaml:"api_port" validate:"min=1,max=65535" default:"7000"`
	EngineConf EngineConf      `yaml:"engine"`
	Schedules  []*ScheduleConf `yaml:"schedules"`
	Mappings   []*MappingConf  `yaml:"mappings"`
	Alerts     []*AlertConf    `yaml:"alerts"`
}

type ScheduleConf struct {
	Query  string
	Period time.Duration
	Delay  time.Duration
	Tags   map[string]string
}

type MappingConf struct {
	Name     string
	Tags     map[string]string
	Detector *DetectorConf
}

type DetectorConf struct {
	NormalCount    int                 `yaml:"normal_count"    validate:"nonzero"       default:"3"`
	Consensus      float64             `yaml:"consensus"       validate:"nonzero"       default:"0.5"`
	MinSupport     float64             `yaml:"min_support"     validate:"nonzero,max=1" default:"0.05"`
	ErrorTolerance float64             `yaml:"error_tolerance" validate:"nonzero,max=1" default:"0.1"`
	Fingerprinters []FingerprinterConf `yaml:"fingerprints"`
}

type AlertConf struct {
	Message   string
	Threshold float64
	Query     string
	Period    time.Duration
	Delay     time.Duration
	Notifiers []*NotifierConf
}

func (self *Config) Default() {
	config.PerformDefault(self)
	self.EngineConf.Default()
	for _, mp := range self.Mappings {
		mp.Default()
	}
}

func (self Config) Validate() (err error) {
	glog.V(2).Info("Validating Config")
	err = validator.Validate(self)
	if err != nil {
		return
	}
	err = self.EngineConf.Validate()
	if err != nil {
		return
	}
	for _, mp := range self.Mappings {
		err = mp.Validate()
		if err != nil {
			return
		}
	}
	return
}

func (self *MappingConf) Validate() error {
	return validator.Validate(self)
}

func (self *MappingConf) Default() {
	config.PerformDefault(self)
	self.Detector.Default()
}

func (self *DetectorConf) Validate() error {
	glog.Info("Validating DetectorConf ", self)
	err := validator.Validate(self)
	if err != nil {
		return err
	}
	for _, f := range self.Fingerprinters {
		err = f.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *DetectorConf) Default() {
	config.PerformDefault(self)
	for _, f := range self.Fingerprinters {
		f.Default()
	}
}

func LoadFromFile(path string) (*Config, error) {
	config_file, err := os.Open(path)
	if err != nil {
		return nil, err

	}
	data, err := ioutil.ReadAll(config_file)
	if err != nil {
		return nil, err
	}
	return Load(data)
}

func Load(data []byte) (*Config, error) {
	config := new(Config)

	err := yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	config.Default()
	err = config.Validate()
	if err != nil {
		return nil, err
	}

	return config, nil
}
