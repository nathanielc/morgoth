package detector

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nathanielc/morgoth/window"
)

type Detector struct {
	normalCount    int
	consensus      float64
	minSupport     float64
	errorTolerance float64
	counters       []fingerprinterCounter
}

// Pair of fingerprinter and counter
type fingerprinterCounter struct {
	fingerprinter Fingerprinter
	counter       Counter
}

func New(normalCount int, consensus, minSupport, errorTolerance float64) *Detector {
	//TODO perform sanity check on minsupport and normalcount to make sure
	// its still possible to mark something as anomalous
	return &Detector{
		normalCount:    normalCount,
		consensus:      consensus,
		minSupport:     minSupport,
		errorTolerance: errorTolerance,
	}
}

// Each fingerprinter will be used in detecting anomalies
func (self *Detector) AddFingerprinter(f Fingerprinter) {
	counter := NewLossyCounter(self.minSupport, self.errorTolerance)
	fc := fingerprinterCounter{
		fingerprinter: f,
		counter:       counter,
	}
	self.counters = append(self.counters, fc)
}

// Determine if the window is anomalous
func (self *Detector) IsAnomalous(window *window.Window) bool {

	vote := 0.0
	for _, fc := range self.counters {
		fingerprint := fc.fingerprinter.Fingerprint(window.Data)
		count := fc.counter.Count(fingerprint)
		glog.Infof("Count: %d", count)
		if count < self.normalCount {
			vote++
		}
	}

	glog.Infof("Anomalous v:%f", vote)

	anomalous := vote/float64(len(self.counters)) >= self.consensus

	return anomalous
}
