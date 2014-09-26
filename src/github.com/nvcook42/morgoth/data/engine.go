package data

type Engine interface {
	GetReader() Reader
	GetWriter() Writer
}
