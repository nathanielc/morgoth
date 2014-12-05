package ymtn_test

import (
	"flag"
	"github.com/nvcook42/morgoth/learners/ymtn"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func init() {
	flag.Parse()
	if testing.Verbose() {
		flag.Set("logtostderr", "true")
		flag.Set("v", "3")
	}
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
	assert.Equal(0.0, sum)
}

func TestRSSTSin(t *testing.T) {
	assert := assert.New(t)

	size := 100
	x := make([]float64, size, size)
	for i := range x {
		x[i] = math.Sin(float64(i+1))
	}
	scores := ymtn.RSST(x, 5, 4)
	assert.NotNil(scores)

	sum := 0.0
	for _, v := range scores {
		sum += v
	}
	assert.Equal(0.0, sum)
}

func TestRSSTSaw(t *testing.T) {
	assert := assert.New(t)

	size := 100
	x := make([]float64, size, size)
	for i := range x {
		x[i] = float64((i+1) % 10)
	}
	scores := ymtn.RSST(x, 5, 4)
	assert.NotNil(scores)

	sum := 0.0
	for _, v := range scores {
		sum += v
	}
	assert.Equal(13.862825859387659, sum)
}
