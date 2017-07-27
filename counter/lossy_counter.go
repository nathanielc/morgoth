package counter

import (
	"math"
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type lossyCounter struct {
	mu                 sync.RWMutex
	errorTolerance     float64
	frequencies        []*entry
	distributionGauges []prometheus.Gauge
	width              int
	total              int
	bucket             int

	metrics *Metrics
}

type entry struct {
	countable Countable
	count     int
	delta     int
}

//Create a new lossycounter with specified errorTolerance
func NewLossyCounter(metrics *Metrics, errorTolerance float64) *lossyCounter {
	return &lossyCounter{
		metrics:        metrics,
		errorTolerance: errorTolerance,
		width:          int(math.Ceil(1.0 / errorTolerance)),
		total:          0,
		bucket:         1,
	}
}

// Count a countable and return the support for the countable.
func (self *lossyCounter) Count(countable Countable) float64 {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.total++

	count := 0
	for i, existing := range self.frequencies {
		if existing.countable.IsMatch(countable) {
			//Found match, count it
			existing.count++
			count = existing.count
			// Keep new countable to allow for drift
			self.frequencies[i].countable = countable

			self.distributionGauges[i].Set(float64(count))
			break
		}
	}

	if count == 0 {
		// No matches create new entry
		count = 1

		// Create new gauge
		lvs := append(self.metrics.LabelValues, strconv.Itoa(len(self.distributionGauges)))
		g := self.metrics.Distribution.WithLabelValues(lvs...)
		g.Set(float64(count))

		// Count new unique fingerprint
		self.metrics.UniqueFingerprints.Inc()

		// append
		self.frequencies = append(self.frequencies, &entry{
			countable: countable,
			count:     count,
			delta:     self.bucket - 1,
		})
		self.distributionGauges = append(self.distributionGauges, g)
	}

	if self.total%self.width == 0 {
		self.prune()
		self.bucket++
	}

	return float64(count) / float64(self.total)
}

//Remove infrequent items from the list
func (self *lossyCounter) prune() {
	filteredFreqs := self.frequencies[0:0]
	filteredGauges := self.distributionGauges[0:0]
	self.metrics.Distribution.Reset()
	for i, entry := range self.frequencies {
		if entry.count+entry.delta > self.bucket {
			lvs := append(self.metrics.LabelValues, strconv.Itoa(i))
			g := self.metrics.Distribution.WithLabelValues(lvs...)
			g.Set(float64(entry.count))

			filteredFreqs = append(filteredFreqs, entry)
			filteredGauges = append(filteredGauges, g)
		}
	}

	self.frequencies = filteredFreqs
	self.distributionGauges = filteredGauges

	self.metrics.UniqueFingerprints.Set(float64(len(self.frequencies)))
}
