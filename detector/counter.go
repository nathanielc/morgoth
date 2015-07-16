package detector

type Counter interface {
	// Count a fingerprint and return the number of times
	// that fingerprint has been seen
	Count(Fingerprint) int
}
