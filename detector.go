package morgoth

import (
	"errors"
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

// Create a new Lossy couting based detector
// The consensus is a percentage of the fingerprinters that must agree in order to flag a window as anomalous.
// If the consensus is -1 then the average support from each fingerprinter is compared to minSupport instead of using a consensus.
// The minSupport defines a minimum frequency as a percentage for a window to be considered normal.
// The errorTolerance defines a frequency as a precentage for the smallest frequency that will be retained in memory.
// The errorTolerance must be less than the minSupport.
func NewDetector(consensus, minSupport, errorTolerance float64, fingerprinters []Fingerprinter) (*Detector, error) {
	if (consensus != -1 && consensus < 0) || consensus > 1 {
		return nil, errors.New("consensus must be in the range [0,1) or equal to -1")
	}
	if minSupport <= errorTolerance {
		return nil, errors.New("minSupport must be greater than errorTolerance")
	}
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
	}, nil
}

// Determine if the window is anomalous
func (self *Detector) IsAnomalous(window *Window) (bool, float64) {
	self.Stats.WindowCount++
	self.Stats.DataPointCount += uint64(len(window.Data))

	vote := 0.0
	avgSupport := 0.0
	n := 0.0
	for _, fc := range self.counters {
		fingerprint := fc.fingerprinter.Fingerprint(window.Copy())
		support := fc.counter.Count(fingerprint)
		anomalous := support <= self.minSupport
		if anomalous {
			vote++
		}
		log.Printf("D! %T anomalous? %v support: %f", fc.fingerprinter, anomalous, support)

		avgSupport = ((avgSupport * n) + support) / (n + 1)
		n++
	}

	anomalous := false
	if self.consensus != -1 {
		// Use voting consensus
		vote /= float64(len(self.counters))
		anomalous = vote >= self.consensus
	} else {
		// Use average suppport
		anomalous = avgSupport <= self.minSupport
	}

	if anomalous {
		self.Stats.AnomalousCount++
	}

	return anomalous, avgSupport
}
