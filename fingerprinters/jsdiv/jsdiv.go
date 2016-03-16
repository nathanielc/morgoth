package jsdiv

import (
	"log"
	"math"

	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/counter"
)

const iterations = 20

var ln2 = math.Log(2)

// Jensen-Shannon Divergence
//
// Fingerprints store the histogram of the window.
// Fingerprints are compared to see their JS divergence distance is less than a critical threshold.
//
// Configuration:
//  min: Minimum value of the window data.
//  max: Maximum value of the window data.
//    NOTE: The JS divergence is symmetrical and I may be able to drop the min/max requirement
//  nBins: Number of bins to use in the histogram.
//  pvalue: Standard p-value statistical threshold. Typical value is 0.05
type JSDiv struct {
	min    float64
	max    float64
	nBins  int
	pValue float64
}

func New(min, max, pValue float64, nBins int) *JSDiv {
	return &JSDiv{
		min:    min,
		max:    max,
		nBins:  nBins,
		pValue: pValue,
	}
}

func (self *JSDiv) Fingerprint(window *morgoth.Window) morgoth.Fingerprint {

	hist, count := calcHistogram(window.Data, self.min, self.max, self.nBins)
	return &JSDivFingerprint{
		hist,
		count,
		self.pValue,
	}
}

func calcHistogram(xs []float64, min, max float64, nBins int) (hist []float64, count int) {
	count = len(xs)
	c := float64(count)
	hist = make([]float64, nBins)
	stepSize := (max - min) / float64(nBins)
	for _, x := range xs {
		i := int(math.Floor((x - min) / stepSize))
		if i > nBins {
			//Just in case x == max value
			i--
		}
		hist[i] += 1 / c
	}
	return
}

type JSDivFingerprint struct {
	histogram []float64
	count     int
	pValue    float64
}

func (self *JSDivFingerprint) IsMatch(other counter.Countable) bool {
	othr, ok := other.(*JSDivFingerprint)
	if !ok {
		return false
	}
	if len(self.histogram) != len(othr.histogram) {
		log.Println("Unexpected comparision between JSDivFingerprints")
		return false
	}

	s := self.calcSignificance(othr)

	return s < self.pValue
}

func (self *JSDivFingerprint) calcSignificance(other *JSDivFingerprint) float64 {
	p := self.histogram
	q := other.histogram
	n := len(p)
	m := make([]float64, n)
	for i := range p {
		m[i] = 0.5 * (p[i] + q[i])
	}

	v := 0.5 * float64(n-1)

	D := calcS(m) - (0.5*calcS(p) + 0.5*calcS(q))

	inc := apporxIncompleteGamma(v, float64(n)*ln2*D)
	gamma := math.Gamma(v)

	return inc / gamma
}

// Calculate the Shannon measure for a histogram
func calcS(hist []float64) float64 {
	s := 0.0
	for _, v := range hist {
		if v != 0 {
			s += v * math.Log2(v)
		}
	}

	return -s
}

// This is a work in progress. Need to update.
func apporxIncompleteGamma(s, x float64) float64 {
	g := 0.0
	xs := math.Pow(x, s)
	ex := math.Exp(-x)

	for k := 0; k < iterations; k++ {
		denominator := s
		for i := 1; i <= k; i++ {
			denominator *= s + float64(i)
		}
		g += (xs * ex * math.Pow(x, float64(k))) / denominator
	}
	return g
}
