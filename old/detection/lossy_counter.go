package detection

import (
	"math"
)

type lossyCounter struct {
	minSupport     float64
	errorTolerance float64
	frequencies    []*entry
	width          int
	total          int
	bucket         int
}

type entry struct {
	fingerprint Fingerprint
	count       int
	delta       int
}

//Create a new lossycounter with specified minSupport and errorTolerance
func NewLossyCounter(minSupport, errorTolerance float64) *lossyCounter {
	return &lossyCounter{
		minSupport:     minSupport,
		errorTolerance: errorTolerance,
		width:          int(math.Ceil(1.0 / errorTolerance)),
		total:          0,
		bucket:         1,
	}
}

// Count a fingerprint and return the number of time that fingerprint has
// been seen within errorTolerance
func (self *lossyCounter) Count(fingerprint Fingerprint) int {

	self.total++

	count := 0
	for _, existing := range self.frequencies {
		if existing.fingerprint.IsMatch(fingerprint) {
			//Found match, count it
			existing.count++
			count = existing.count
			break
		}
	}

	if count == 0 {
		// No matches create new entry
		count = 1
		self.frequencies = append(self.frequencies, &entry{
			fingerprint: fingerprint,
			count:       count,
			delta:       self.bucket - 1,
		})
	}

	if self.total%self.width == 0 {
		self.prune()
		self.bucket++
	}

	return count
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
}
