package morgoth

import (
	"log"
	"sync"

	"github.com/nathanielc/morgoth/counter"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

type Detector struct {
	mu             sync.RWMutex
	consensus      float64
	minSupport     float64
	errorTolerance float64
	counters       []fingerprinterCounter

	metrics *DetectorMetrics
}

type DetectorMetrics struct {
	WindowCount          prometheus.Counter
	PointCount           prometheus.Counter
	AnomalousCount       prometheus.Counter
	FingerprinterMetrics []*counter.Metrics
}

func (m *DetectorMetrics) Register() error {
	if err := prometheus.Register(m.WindowCount); err != nil {
		return errors.Wrap(err, "window count metric")
	}
	if err := prometheus.Register(m.PointCount); err != nil {
		return errors.Wrap(err, "point count metric")
	}
	if err := prometheus.Register(m.AnomalousCount); err != nil {
		return errors.Wrap(err, "anomalous count metric")
	}
	for i, f := range m.FingerprinterMetrics {
		if err := f.Register(); err != nil {
			return errors.Wrapf(err, "fingerprinter %d", i)
		}
	}
	return nil
}
func (m *DetectorMetrics) Unregister() {
	prometheus.Unregister(m.WindowCount)
	prometheus.Unregister(m.PointCount)
	prometheus.Unregister(m.AnomalousCount)
	for _, f := range m.FingerprinterMetrics {
		f.Unregister()
	}
}

// Pair of fingerprinter and counter
type fingerprinterCounter struct {
	Fingerprinter
	counter.Counter
}

// Create a new Lossy couting based detector
// The consensus is a percentage of the fingerprinters that must agree in order to flag a window as anomalous.
// If the consensus is -1 then the average support from each fingerprinter is compared to minSupport instead of using a consensus.
// The minSupport defines a minimum frequency as a percentage for a window to be considered normal.
// The errorTolerance defines a frequency as a precentage for the smallest frequency that will be retained in memory.
// The errorTolerance must be less than the minSupport.
func NewDetector(metrics *DetectorMetrics, consensus, minSupport, errorTolerance float64, fingerprinters []Fingerprinter) (*Detector, error) {
	if (consensus != -1 && consensus < 0) || consensus > 1 {
		return nil, errors.New("consensus must be in the range [0,1) or equal to -1")
	}
	if minSupport <= errorTolerance {
		return nil, errors.New("minSupport must be greater than errorTolerance")
	}
	if len(metrics.FingerprinterMetrics) != len(fingerprinters) {
		return nil, errors.New("must provide the same number of fingerprinter metrics as fingerprinters")
	}
	counters := make([]fingerprinterCounter, len(fingerprinters))
	for i, fingerprinter := range fingerprinters {
		counters[i] = fingerprinterCounter{
			Fingerprinter: fingerprinter,
			Counter:       counter.NewLossyCounter(metrics.FingerprinterMetrics[i], errorTolerance),
		}
	}
	return &Detector{
		metrics:        metrics,
		consensus:      consensus,
		minSupport:     minSupport,
		errorTolerance: errorTolerance,
		counters:       counters,
	}, nil
}

// Determine if the window is anomalous
func (self *Detector) IsAnomalous(window *Window) (bool, float64) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.metrics.WindowCount.Inc()
	self.metrics.PointCount.Add(float64(len(window.Data)))

	vote := 0.0
	avgSupport := 0.0
	n := 0.0
	for _, fc := range self.counters {
		fingerprint := fc.Fingerprint(window.Copy())
		support := fc.Count(fingerprint)
		anomalous := support <= self.minSupport
		if anomalous {
			vote++
		}
		log.Printf("D! %T anomalous? %v support: %f", fc.Fingerprinter, anomalous, support)

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
		self.metrics.AnomalousCount.Inc()
	}

	return anomalous, avgSupport
}

func (self *Detector) Close() {
	self.metrics.Unregister()
}
