package morgoth

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nathanielc/morgoth/counter"
)

type Detector struct {
	normalCount    int
	consensus      float64
	minSupport     float64
	errorTolerance float64
	counters       []fingerprinterCounter
	Stats          DetectorStats
}

type DetectorStats struct {
	WindowCount    uint64
	DataPointCount uint64
	AnomalousCount uint64
}

type DetectorBuilder func() *Detector

// Pair of fingerprinter and counter
type fingerprinterCounter struct {
	fingerprinter Fingerprinter
	counter       counter.Counter
}

func NewDetectorBuilder(normalCount int, consensus, minSupport, errorTolerance float64, fingerprinters []Fingerprinter) DetectorBuilder {
	return func() *Detector {
		return NewDetector(normalCount, consensus, minSupport, errorTolerance, fingerprinters)
	}
}

func NewDetector(normalCount int, consensus, minSupport, errorTolerance float64, fingerprinters []Fingerprinter) *Detector {
	//TODO perform sanity check on minsupport and normalcount to make sure
	// its still possible to mark something as anomalous
	counters := make([]fingerprinterCounter, len(fingerprinters))
	for i, fingerprinter := range fingerprinters {
		counters[i] = fingerprinterCounter{
			fingerprinter,
			counter.NewLossyCounter(minSupport, errorTolerance),
		}
	}
	return &Detector{
		normalCount:    normalCount,
		consensus:      consensus,
		minSupport:     minSupport,
		errorTolerance: errorTolerance,
		counters:       counters,
	}
}

// Determine if the window is anomalous
func (self *Detector) IsAnomalous(window *Window) bool {

	self.Stats.WindowCount++
	self.Stats.DataPointCount += uint64(len(window.Data))

	vote := 0.0
	for _, fc := range self.counters {
		fingerprint := fc.fingerprinter.Fingerprint(window.Copy())
		count := fc.counter.Count(fingerprint)
		if count < self.normalCount {
			vote++
		}
	}

	vote /= float64(len(self.counters))
	anomalous := vote >= self.consensus

	glog.V(3).Infof("Window anomalous: %v vote: %f, consensus: %f", anomalous, vote, self.consensus)

	if anomalous {
		self.Stats.AnomalousCount++
	}

	return anomalous
}
