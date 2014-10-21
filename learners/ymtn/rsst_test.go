package ymtn_test

import (
	"github.com/nvcook42/morgoth/learners/ymtn"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestRSST(t *testing.T) {
	assert := assert.New(t)

	size := 100
	rand.Seed(42)
	x := make([]float64, size, size)
	for i := range x {
		x[i] = rand.Float64()
	}
	xs := ymtn.RSST(x, 5, 4)
	assert.NotNil(xs)

}
