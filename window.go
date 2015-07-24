package morgoth

import (
	"time"
)

type Window struct {
	Name  string
	Data  []float64
	Tags  map[string]string
	Start time.Time
	Stop  time.Time
}

// Search all given tags to see if they match the window's tags
func (self *Window) IsTagsMatch(tags map[string]string) bool {
	for k, v := range tags {
		if tag, ok := self.Tags[k]; !ok || tag != v {
			return false
		}
	}
	return true
}

func (self *Window) Copy() Window {

	data := make([]float64, len(self.Data))
	copy(data, self.Data)

	tags := make(map[string]string, len(self.Tags))
	for k, v := range self.Tags {
		tags[k] = v
	}

	return Window{
		Name:  self.Name,
		Data:  data,
		Tags:  tags,
		Start: self.Start,
		Stop:  self.Stop,
	}
}
