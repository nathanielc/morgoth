package detection

type Fingerprint interface {
	IsMatch(other Fingerprint) bool
}
