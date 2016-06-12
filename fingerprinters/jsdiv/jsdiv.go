package jsdiv

import (
	"math"

	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/counter"
)

const iterations = 20

type histogram map[int]float64

var ln2 = math.Log(2)

// Jensen-Shannon Divergence
//
// Fingerprints store the histogram of the window.
// Fingerprints are compared to see their JS divergence distance is less than a critical threshold.
//
// Configuration:
//  min: Excpected minimum value of the window data.
//  max: Excpected maximum value of the window data.
//  binwidth: Size of a bin for the histogram
//  pValue: Standard p-value statistical threshold. Typical value is 0.05
type JSDiv struct {
	minIndex int
	maxIndex int
	binWidth float64
	pValue   float64
}

func New(min, max, binWidth, pValue float64) *JSDiv {
	return &JSDiv{
		minIndex: int(math.Floor(min / binWidth)),
		maxIndex: int(math.Floor(max / binWidth)),
		binWidth: binWidth,
		pValue:   pValue,
	}
}

func (self *JSDiv) Fingerprint(window *morgoth.Window) morgoth.Fingerprint {

	hist, count := calcHistogram(window.Data, self.binWidth)
	return &JSDivFingerprint{
		hist,
		count,
		self.pValue,
		self.minIndex,
		self.maxIndex,
	}
}

func calcHistogram(xs []float64, binWidth float64) (hist histogram, count int) {
	count = len(xs)
	c := float64(count)
	hist = make(histogram)
	for _, x := range xs {
		i := int(math.Floor(x / binWidth))
		hist[i] += 1.0 / c
	}
	return
}

type JSDivFingerprint struct {
	histogram histogram
	count     int
	pValue    float64

	minIndex int
	maxIndex int
}

func (self *JSDivFingerprint) IsMatch(other counter.Countable) bool {
	othr, ok := other.(*JSDivFingerprint)
	if !ok {
		return false
	}

	s := self.calcSignificance(othr)

	return s < self.pValue
}

func (self *JSDivFingerprint) calcSignificance(other *JSDivFingerprint) float64 {
	p := self.histogram
	q := other.histogram
	m := make(histogram, len(p)+len(q))
	min := self.minIndex
	max := self.maxIndex
	for i := range p {
		if i < min {
			min = i
		}
		if i > max {
			max = i
		}
		m[i] = 0.5 * p[i]
	}
	for i := range q {
		if i < min {
			min = i
		}
		if i > max {
			max = i
		}
		m[i] += 0.5 * q[i]
	}

	k := max - min

	v := 0.5 * float64(k-1)

	D := calcS(m) - (0.5*calcS(p) + 0.5*calcS(q))

	inc := apporxIncompleteGamma(v, float64(self.count+other.count)*ln2*D)
	gamma := math.Gamma(v)

	return inc / gamma
}

// Calculate the Shannon measure for a histogram
func calcS(hist histogram) float64 {
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
