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
	Distribution       *prometheus.GaugeVec
}

func (m *Metrics) Register() error {
	if err := prometheus.Register(m.UniqueFingerprints); err != nil {
		return errors.Wrap(err, "unique fingerprints metric")
	}
	if err := prometheus.Register(m.Distribution); err != nil {
		return errors.Wrap(err, "distribution metric")
	}
	return nil
}

func (m *Metrics) Unregister() {
	prometheus.Unregister(m.UniqueFingerprints)
	prometheus.Unregister(m.Distribution)
}
