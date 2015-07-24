package counter

type Counter interface {
	// Count a fingerprint and return the number of times
	// that fingerprint has been seen
	Count(Countable) int
}

type Countable interface {
	IsMatch(other Countable) bool
}
