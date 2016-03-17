package counter

type Counter interface {
	// Count a fingerprint and return the support for the item.
	// support = count / total
	Count(Countable) float64
}

type Countable interface {
	IsMatch(other Countable) bool
}
