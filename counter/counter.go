package counter

import (
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

type Counter interface {
	// Count a fingerprint and return the support for the item.
	// support = count / total
	Count(Countable) float64
}

type Countable interface {
	IsMatch(other Countable) bool
}

type Metrics struct {
	UniqueFingerprints prometheus.Gauge
	// TODO(nathanielc): Figure out how to represent
	// this distribution data as a prometheus.
	//Distribution []int
}

func (m *Metrics) Register() error {
	return errors.Wrap(prometheus.Register(m.UniqueFingerprints), "unique fingerprints metric")
}

func (m *Metrics) Unregister() {
	prometheus.Unregister(m.UniqueFingerprints)
}
