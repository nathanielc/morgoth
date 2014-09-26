package config

// Pattern is regexp
type Pattern string

// Represents a single metric conf section
type Metric struct {
	Pattern   Pattern    `yaml:"pattern"`
	Schedule  Schedule   `yaml:"schedule"`
	Detectors []Detector `yaml:"detectors"`
}

func (self *Metric) Default() {
	self.Schedule.Default()

	for i := range self.Detectors {
		self.Detectors[i].Default()
	}
}

func (self Metric) Validate() error {
	if valid := self.Pattern.Validate(); valid != nil {
		return valid
	}
	if valid := self.Schedule.Validate(); valid != nil {
		return valid
	}
	for i := range self.Detectors {
		if valid := self.Detectors[i].Validate(); valid != nil {
			return valid
		}
	}
	return nil
}
