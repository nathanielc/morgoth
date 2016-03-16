package morgoth

import "github.com/nathanielc/morgoth/counter"

type Fingerprint interface {
	IsMatch(other counter.Countable) bool
}

type Fingerprinter interface {
	Fingerprint(window *Window) Fingerprint
}
