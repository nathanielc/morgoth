package morgoth

import (
	"hash/fnv"
	"regexp"
	"sort"
)

type Mapper struct {
	detectorMappers  map[uint64][]*DetectorMapper
	detectorMatchers []*DetectorMatcher
	Stats            MapperStats
}

type MapperStats struct {
	WindowCount  uint64
	MapperStats  []*DetectorMapperStats
	MatcherStats []*DetectorMatcherStats
}

type DetectorMatcher struct {
	NamePattern     *regexp.Regexp
	TagPatterns     map[string]*regexp.Regexp
	DetectorBuilder func() *Detector
	Stats           DetectorMatcherStats
}

type DetectorMatcherStats struct {
	NamePattern string
	TagPatterns map[string]string
	MatchCount  uint64
}

type DetectorMapper struct {
	Name     string
	Tags     map[string]string
	detector *Detector
	Stats    DetectorMapperStats
}

type DetectorMapperStats struct {
	Name          string
	Tags          map[string]string
	MapCount      uint64
	DetectorStats *DetectorStats
}

func NewMapper(detectorMappers []*DetectorMapper, detectorMatchers []*DetectorMatcher) *Mapper {
	mapper := &Mapper{
		detectorMappers: make(
			map[uint64][]*DetectorMapper,
			len(detectorMappers),
		),
		detectorMatchers: make(
			[]*DetectorMatcher,
			0,
			len(detectorMatchers),
		),
	}
	for _, m := range detectorMappers {
		mapper.addDetectorMapper(m)
	}
	for _, m := range detectorMatchers {
		mapper.addDetectorMatcher(m)
	}
	return mapper
}

func (self *Mapper) addDetectorMapper(mapper *DetectorMapper) {
	hash := calcHash(mapper.Name, mapper.Tags)
	self.detectorMappers[hash] = append(
		self.detectorMappers[hash],
		mapper,
	)
	self.Stats.MapperStats = append(
		self.Stats.MapperStats,
		&mapper.Stats,
	)
}

func (self *Mapper) addDetectorMatcher(matcher *DetectorMatcher) {
	self.detectorMatchers = append(
		self.detectorMatchers,
		matcher,
	)
	self.Stats.MatcherStats = append(
		self.Stats.MatcherStats,
		&matcher.Stats,
	)
}

func (self *Mapper) Map(w *Window) *Detector {

	self.Stats.WindowCount++

	// First check for match in detector mappers
	hash := calcHash(w.Name, w.Tags)
	for _, mapper := range self.detectorMappers[hash] {
		detector := mapper.Map(w)
		if detector != nil {
			return detector
		}
	}

	// Last check for match in Detector matchers
	for _, matcher := range self.detectorMatchers {
		if matcher.IsMatch(w) {
			mapper := matcher.NewDetectorMapper(w)
			self.addDetectorMapper(mapper)
			return mapper.detector
		}
	}
	// No mapping found
	return nil
}

func (self *Mapper) GetDetectorMappings() []*DetectorMapper {
	mappings := make([]*DetectorMapper, 0, len(self.detectorMappers))
	for _, value := range self.detectorMappers {
		mappings = append(mappings, value...)
	}
	return mappings
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
	return hash.Sum64()
}

func NewDetectorMapper(w *Window, detector *Detector) *DetectorMapper {
	m := &DetectorMapper{
		Name:     w.Name,
		Tags:     w.Tags,
		detector: detector,
		Stats: DetectorMapperStats{
			Name:          w.Name,
			Tags:          w.Tags,
			DetectorStats: &detector.Stats,
		},
	}
	return m
}

func (self *DetectorMapper) Map(w *Window) *Detector {
	if self.Name != w.Name {
		return nil
	}

	// Check that window doesn't have extra tags
	if len(self.Tags) != len(w.Tags) {
		return nil
	}

	// Check that tag sets match
	for k, mapTag := range self.Tags {
		if tag, ok := w.Tags[k]; !ok || tag != mapTag {
			return nil
		}
	}

	self.Stats.MapCount++

	return self.detector
}

func NewDetectorMatcher(namePattern *regexp.Regexp, tagPatterns map[string]*regexp.Regexp, detectorBuilder DetectorBuilder) *DetectorMatcher {
	tags := make(map[string]string, len(tagPatterns))
	for tag, pattern := range tagPatterns {
		tags[tag] = pattern.String()
	}
	return &DetectorMatcher{
		NamePattern:     namePattern,
		TagPatterns:     tagPatterns,
		DetectorBuilder: detectorBuilder,
		Stats: DetectorMatcherStats{
			NamePattern: namePattern.String(),
			TagPatterns: tags,
		},
	}
}

func (self *DetectorMatcher) IsMatch(w *Window) bool {
	if !self.NamePattern.MatchString(w.Name) {
		return false
	}

	//Check only defined tags match patterns
	for k, pattern := range self.TagPatterns {
		if tag, ok := w.Tags[k]; !ok || !pattern.MatchString(tag) {
			return false
		}
	}

	self.Stats.MatchCount++

	return true
}

func (self *DetectorMatcher) NewDetectorMapper(w *Window) *DetectorMapper {
	detector := self.DetectorBuilder()
	return NewDetectorMapper(w, detector)
}
