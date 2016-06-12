package kstest

import (
	"math"
	"sort"

	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/counter"
)

var confidenceMappings = []float64{
	1.22,
	1.36,
	1.48,
	1.63,
	1.73,
	1.95,
}

// Kolmogorov–Smirnov test.
// https://en.wikipedia.org/wiki/Kolmogorov%E2%80%93Smirnov_test
//
// The fingerprint is the cummulative distribution of the window.
// The fingerprints are compared by computing the largest distance between the cummulative distribution functions and comparing to a critical value.
//
// Configuration:
//  The only parameter is a confidence level.
//  Valid values are from 0-5.
//  The level maps to a list of predefined critical values for the KS test.
//  Increasing 'confidence' decreases the number of anomalies detected.
//
type KSTest struct {
	confidence uint
}

func New(confidence uint) *KSTest {
	return &KSTest{
		confidence: confidence,
	}
}

func (self *KSTest) Fingerprint(window *morgoth.Window) morgoth.Fingerprint {

	sort.Float64s(window.Data)

	return &KSTestFingerprint{self.confidence, window.Data}
}

type KSTestFingerprint struct {
	confidence uint
	edf        []float64
}

func (self *KSTestFingerprint) IsMatch(other counter.Countable) bool {
	othr, ok := other.(*KSTestFingerprint)
	if !ok {
		return false
	}
	if self.confidence != othr.confidence {
		return false
	}

	threshold := self.calcThreshold(othr)

	D := calcD(self.edf, othr.edf)

	return D < threshold
}

// Calculate the critical threshold for this comparision
func (self *KSTestFingerprint) calcThreshold(othr *KSTestFingerprint) float64 {
	c := confidenceMappings[self.confidence]
	n := float64(len(self.edf))
	m := float64(len(othr.edf))
	return c * math.Sqrt((n+m)/(n*m))

}

// Calculate maximum distance between cummulative distributions
func calcD(f1, f2 []float64) float64 {
	D := 0.0
	n := len(f1)
	m := len(f2)
	i := 0
	j := 0
	for i < n && j < m {
		for i < n && j < m && f1[i] < f2[j] {
			i++
		}
		for i < n && j < m && f1[i] > f2[j] {
			j++
		}
		for i < n && j < m && f1[i] == f2[j] {
			i++
			j++
		}
		cdf1 := float64(i) / float64(n)
		cdf2 := float64(j) / float64(m)
		if d := math.Abs(cdf1 - cdf2); d > D {
			D = d
		}
	}
	return D
}
