package counter

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

//Simple fingerprint implementation
type fp struct {
	id int
}

func (self *fp) IsMatch(other Countable) bool {
	fp, ok := other.(*fp)
	return ok && self.id == fp.id
}

func TestLossyCounterShouldCountAllItems(t *testing.T) {
	assert := assert.New(t)

	lc := NewLossyCounter(0.01)

	fp1 := &fp{1}
	fp2 := &fp{2}

	assert.NotEqual(fp1, fp2)

	assert.Equal(1.0/1.0, lc.Count(fp1))
	assert.Equal(2.0/2.0, lc.Count(fp1))
	assert.Equal(1.0/3.0, lc.Count(fp2))
	assert.Equal(2.0/4.0, lc.Count(fp2))
	assert.Equal(3.0/5.0, lc.Count(fp1))
	assert.Equal(4.0/6.0, lc.Count(fp1))
}

func TestLossyCounterShouldByLossy(t *testing.T) {
	assert := assert.New(t)

	//Create Lossy Counter that will drop items less than 10% frequent
	lc := NewLossyCounter(0.10)

	fp1 := &fp{1}
	fp2 := &fp{2}

	// Count fp1 10 times: 10%
	for i := 0; i < 10; i++ {
		assert.Equal(1.0, lc.Count(fp1))
	}

	// Count fp2 90 times: 90%
	for i := 0; i < 90; i++ {
		assert.Equal(float64(i+1)/float64(11+i), lc.Count(fp2))
	}

	// Count fp1 once more, should have been dropped and
	// now is counted only once
	assert.Equal(1.0/101.0, lc.Count(fp1))
}

//Benchmark the worst case scenario for the lossy counter:
// every item is errorTolerance frequent
func BenchmarkCounting(b *testing.B) {

	e := 0.01
	lc := NewLossyCounter(0.01)

	unique := int(math.Ceil(1.0 / e))

	fps := make([]*fp, unique)
	for i := 0; i < unique; i++ {
		fps[i] = &fp{i}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := i % unique
		lc.Count(fps[id])
	}

}
