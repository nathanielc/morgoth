package config

// A Factory provides methods to create an empty configuration object
// and to get an instance from a populated configuration
type Factory interface {
	//Return a new zeroed configuration.
	NewConf() Configuration
	//Return a new instance based on the given configuration.
	GetInstance(Configuration) (interface{}, error)
}
