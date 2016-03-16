package sigma

import (
	"math"

	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/counter"
)

// Simple fingerprinter that computes both mean and standard deviation of a window.
// Fingerprints are compared to see if the means are more than n deviations apart.
type Sigma struct {
	deviations float64
}

func New(deviations float64) *Sigma {
	return &Sigma{
		deviations: deviations,
	}
}

func (self *Sigma) Fingerprint(window *morgoth.Window) morgoth.Fingerprint {
	mean, std := calcStats(window.Data)
	return SigmaFingerprint{
		mean:      mean,
		threshold: self.deviations * std,
	}
}

func calcStats(xs []float64) (mean, std float64) {
	n := 0.0
	M2 := 0.0

	for _, x := range xs {
		n++
		delta := x - mean
		mean = mean + delta/n
		M2 += delta * (x - mean)
	}

	std = math.Sqrt(M2 / n)
	return
}

type SigmaFingerprint struct {
	mean      float64
	threshold float64
}

func (self SigmaFingerprint) IsMatch(other counter.Countable) bool {
	o, ok := other.(SigmaFingerprint)
	if !ok {
		return false
	}
	return math.Abs(self.mean-o.mean) < self.threshold
}
