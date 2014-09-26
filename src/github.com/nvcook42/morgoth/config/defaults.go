package config

func (self *Config) Default() {
	self.DataEngine.Default()
	self.Fittings.Default()
}

func (self *DataEngine) Default() {
	switch self.Type {
	case InfluxDB:
		self.InfluxDB.Default()
	case MongoDB:
		self.MongoDB.Default()
	}
}

func (self *Fittings) Default() {
}
