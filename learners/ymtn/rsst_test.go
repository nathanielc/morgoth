package ymtn_test

import (
	"flag"
	"github.com/nvcook42/morgoth/learners/ymtn"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	flag.Set("logtostderr", "true")
	flag.Set("v", "1")
}

func TestRSSTTrivialCase(t *testing.T) {
	assert := assert.New(t)

	size := 100
	x := make([]float64, size, size)
	for i := range x {
		x[i] = 1
	}
	scores := ymtn.RSST(x, 5, 4)
	assert.NotNil(scores)

	sum := 0.0
	for _, v := range scores {
		sum += v
	}
	assert.Equal(2.0, sum)
}
