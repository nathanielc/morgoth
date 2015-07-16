package morgoth

import (
	"hash/fnv"
	"regexp"
)

type Mapper struct {
	detectorMaps     map[uint64]*DetectorMap
	detectorConfMaps *DetectorConfMap
}

func (self *Mapper) Map(w *Window) *Detector {

	// First check for match in detector maps
	hash := calcHash()
	for _, detectorMap := range self.detectorMaps {
		if detectorMap.IsMatch(w) {
			return detectorMap.GetDetector()
		}
	}

	// Last check for match in dectector conf maps
	for _, detectorConfMap := range self.detectorConfMap {
		if detectorConfMap.IsMatch(w) {
			detectorMap := detectorConfMap.NewDetectorMap(w)
			self.detectorMaps = append(self.detectorMaps, detectorMap)
			return detectorMap.GetDetector()
		}
	}
	// No mapping found
	return nil
}

func calcHash(name string, tags map[string]string) uint64 {
	hash := fnv.New64a()
	hash.Write([]byte(name))

	// Sort all tag keys
	keys := make([]string, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	//Hash key/value pairs
	for _, k := range keys {
		hash.Write([]byte(k))
		hash.Write([]byte(tags[k]))
	}
	return fnv.Sum64()
}

type DetectorMap struct {
	Name     string
	Tags     map[string]string
	detector *Detector
}

func NewDetectorMap(w *Window, detector *Detector) *DetectorMap {
	return &DetectorMap{
		Name:     w.Name,
		Tags:     w.Tags,
		detector: detector,
	}
}

func (self *DetectorMap) IsMatch(w *window) bool {
	if self.Name != w.Name {
		return false
	}

	// Check that window doesn't have extra tags
	if len(self.Tags) != len(w.Tags) {
		return false
	}

	// Check that tag sets match
	for k, mapTag := range self.Tags {
		if tag, ok := w.Tags[k]; !ok || tag != mapTag {
			return false
		}
	}

	return true
}

func (self *DetectorMap) GetDetector() *Detector {
	if self.detector != nil {
		return self.detector
	}

	// Go fetch detector from db
}

type DetectorConfMap struct {
	NamePattern         *regexp.Regexp
	TagPatterns         map[string]*regexp.Regexp
	detectorConstructor func() *Detector
}

func (self *DetectorConfMap) IsMatch(w *Window) bool {
	if !self.NamePattern.MatchString(w.Name) {
		return false
	}

	//Check only defined tags match patterns
	for k, pattern := range self.TagPatterns {
		if tag, ok := w.Tags[k]; !ok || !pattern.MatchString(tag) {
			return false
		}
	}

	return true
}

func (self *DetectorConfMap) NewDetectorMap(w *Window) *DetectorMap {
	return NewDetectorMap(w, self.detectorConstructor())
}
