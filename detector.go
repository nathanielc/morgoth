package morgoth

import (
	"log"

	"github.com/nathanielc/morgoth/counter"
)

type Detector struct {
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

// Pair of fingerprinter and counter
type fingerprinterCounter struct {
	fingerprinter Fingerprinter
	counter       counter.Counter
}

func NewDetector(consensus, minSupport, errorTolerance float64, fingerprinters []Fingerprinter) *Detector {
	counters := make([]fingerprinterCounter, len(fingerprinters))
	for i, fingerprinter := range fingerprinters {
		counters[i] = fingerprinterCounter{
			fingerprinter,
			counter.NewLossyCounter(errorTolerance),
		}
	}
	return &Detector{
		consensus:      consensus,
		minSupport:     minSupport,
		errorTolerance: errorTolerance,
		counters:       counters,
	}
}

// Determine if the window is anomalous
func (self *Detector) IsAnomalous(window *Window) (bool, float64) {

	self.Stats.WindowCount++
	self.Stats.DataPointCount += uint64(len(window.Data))

	vote := 0.0
	for _, fc := range self.counters {
		fingerprint := fc.fingerprinter.Fingerprint(window.Copy())
		support := fc.counter.Count(fingerprint)
		log.Printf("F: %T anomalous? %v support: %f", fc.fingerprinter, support < self.minSupport, support)
		if support < self.minSupport {
			vote++
		}
	}

	vote /= float64(len(self.counters))
	anomalous := vote >= self.consensus

	if anomalous {
		self.Stats.AnomalousCount++
	}

	return anomalous, vote
}
