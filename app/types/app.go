// Define exported types from `app' to avoid cyclic imports
package types

import (
	"github.com/nvcook42/morgoth/engine"
	"github.com/nvcook42/morgoth/schedule"
)

// An App is the center on morgoth connecting all the various components
type App interface {
	Run() error
	GetWriter() engine.Writer
	GetReader() engine.Reader
	GetSchedule() schedule.Schedule
}
