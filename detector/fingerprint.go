package detector

type Fingerprint interface {
	IsMatch(other Fingerprint) bool
}
