package morgoth

import (
	"github.com/nathanielc/morgoth/counter"
)

type Fingerprinter interface {
	Fingerprint(window []float64) Fingerprint
}

type Fingerprint interface {
	IsMatch(other counter.Countable) bool
}
