package router

import (
	"github.com/nathanielc/morgoth/detection"
	"github.com/nathanielc/morgoth/window"
	"regexp"
)

type Router interface {
	Route(window *window.Window) bool
}

type UnicastRouter struct {
	routes []Router
}

func (self *UnicastRouter) Route(window *window.Window) {
	for _, route := range self.routes {
		if route.Route(window) {
			break
		}
	}
}

type FanoutRouter struct {
	NamePattern         *regexp.Regexp
	TagPatterns         map[string]*regexp.Regexp
	detectors           []*DetectorRoute
	detectorConstructor func() *detection.Detection
}

func (self *FanoutRouter) Route(window *window.Window) bool {
	if !self.IsMatch(window) {
		return false
	}

	for _, detectorRoute := range self.detectors {
		if detectorRoute.Route(window) {
			return true
		}
	}

	self.detectors = append(
		self.detectors,
		&DetectorRoute{
			Tags:     window.Tags,
			detector: self.detectorConstructor(),
		},
	)

	return true
}

func (self *FanoutRouter) IsMatch(window *window.Window) bool {
	if !self.NamePattern.MatchString(window.Name) {
		return false
	}

	for k, pattern := range self.TagPatterns {
		if tag, ok := window.Tags[k]; !ok || !pattern.MatchString(tag) {
			return false
		}
	}

	return true
}

type DetectorRoute struct {
	Tags     map[string]string
	detector *detection.Detection
}

func (self *DetectorRoute) Route(window *window.Window) bool {
	if !self.IsMatch(window) {
		return false
	}

	self.detector.IsAnomalous(window)
	return true
}

func (self *DetectorRoute) IsMatch(window *window.Window) bool {
	for k, v := range self.Tags {
		if tag, ok := window.Tags[k]; !ok || tag != v {
			return false
		}
	}
	return true
}
