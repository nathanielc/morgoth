// Define exported types from `app' to avoid cyclic imports
package types

import (
	"github.com/nathanielc/morgoth/detector/metadata"
	"github.com/nathanielc/morgoth/engine"
	"github.com/nathanielc/morgoth/schedule"
)

// An App is the center on morgoth connecting all the various components
type App interface {
	Run() error
	GetWriter() engine.Writer
	GetReader() engine.Reader
	GetSchedule() schedule.Schedule
	GetMetadataStore(detectorID string) (metadata.MetadataStore, error)
}
