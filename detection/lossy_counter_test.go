package detection_test

import (
	"flag"
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nvcook42/morgoth/detection"
	"math"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	defer glog.Flush()
	flag.Parse()
	if testing.Verbose() {
		flag.Set("logtostderr", "1")
	}
	os.Exit(m.Run())
}

//Simple fingerprint implementation
type fp struct {
	id int
}

func (self *fp) IsMatch(other detection.Fingerprint) bool {
	return self == other
}

func TestLossyCounterShouldCountAllItems(t *testing.T) {
	assert := assert.New(t)

	lc := detection.NewLossyCounter(0.05, 0.01)

	fp1 := &fp{1}
	fp2 := &fp{2}

	assert.NotEqual(fp1, fp2)

	assert.Equal(1, lc.Count(fp1))
	assert.Equal(2, lc.Count(fp1))
	assert.Equal(1, lc.Count(fp2))
	assert.Equal(2, lc.Count(fp2))
	assert.Equal(3, lc.Count(fp1))
	assert.Equal(4, lc.Count(fp1))
}

func TestLossyCounterShouldByLossy(t *testing.T) {
	assert := assert.New(t)

	//Create Lossy Counter that will drop items less than 10% frequent
	lc := detection.NewLossyCounter(0.5, 0.1)

	fp1 := &fp{1}
	fp2 := &fp{2}

	// Count fp1 10 times: 10%
	for i := 0; i < 10; i++ {
		assert.Equal(i+1, lc.Count(fp1))
	}

	// Count fp2 90 times: 90%
	for i := 0; i < 90; i++ {
		assert.Equal(i+1, lc.Count(fp2))
	}

	// Count fp1 once more, should have been dropped and
	// now is counted only once
	assert.Equal(1, lc.Count(fp1))
}

//Benchmark the worst case scenario for the lossy counter:
// every item is errorTolerance frequent
func BenchmarkCounting(b *testing.B) {

	e := 0.01
	lc := detection.NewLossyCounter(0.5, e)
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
