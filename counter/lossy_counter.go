package counter

import (
	"math"
	"sync"
)

type lossyCounter struct {
	mu             sync.RWMutex
	errorTolerance float64
	frequencies    []*entry
	width          int
	total          int
	bucket         int

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
			break
		}
	}

	if count == 0 {
		// No matches create new entry
		count = 1
		self.frequencies = append(self.frequencies, &entry{
			countable: countable,
			count:     count,
			delta:     self.bucket - 1,
		})
		self.metrics.UniqueFingerprints.Inc()
	}

	if self.total%self.width == 0 {
		self.prune()
		self.bucket++
	}

	return float64(count) / float64(self.total)
}

//Remove infrequent items from the list
func (self *lossyCounter) prune() {
	i := 0
	for i < len(self.frequencies) {
		entry := self.frequencies[i]
		if entry.count+entry.delta <= self.bucket {
			self.frequencies = append(self.frequencies[:i], self.frequencies[i+1:]...)
		} else {
			i++
		}
	}
	self.metrics.UniqueFingerprints.Set(float64(len(self.frequencies)))
}
