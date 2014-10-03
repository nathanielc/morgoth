package registery

type Factory interface {
	NewConf() Configuration
	GetInstance(Configuration) (interface{}, error)
}
