package detection

type Fingerprinter interface {
	Fingerprint(window []float64) Fingerprint
}
