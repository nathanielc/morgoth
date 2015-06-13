package window

type Window struct {
	Name string
	Data []float64
	Tags map[string]string
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
