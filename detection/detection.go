package detection

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/golang/glog"
)

type Detection struct {
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

func New(normalCount int, consensus, minSupport, errorTolerance float64) *Detection {
	//TODO perform sanity check on minsupport and normalcount to make sure
	// its still possible to mark something as anomalous
	return &Detection{
		normalCount:    normalCount,
		consensus:      consensus,
		minSupport:     minSupport,
		errorTolerance: errorTolerance,
	}
}

// Add a fingerprinter to this detection
// Each fingerprinter will be used in detecting anomalies
func (self *Detection) AddFingerprinter(f Fingerprinter) {
	counter := NewLossyCounter(self.minSupport, self.errorTolerance)
	fc := fingerprinterCounter{
		fingerprinter: f,
		counter:       counter,
	}
	self.counters = append(self.counters, fc)
}

// Determine if the window is anomalous
func (self *Detection) IsAnomalous(window []float64) bool {

	vote := 0.0
	for _, fc := range self.counters {
		fingerprint := fc.fingerprinter.Fingerprint(window)
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
