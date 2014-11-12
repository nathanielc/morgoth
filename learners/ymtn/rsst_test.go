package ymtn_test

import (
	"github.com/nvcook42/morgoth/learners/ymtn"
	"github.com/stretchr/testify/assert"
	"math/rand"
	//"math"
	"testing"
)

func TestRSST(t *testing.T) {
	assert := assert.New(t)

	size := 100
	rand.Seed(42)
	x := make([]float64, size, size)
	for i := range x {
		//x[i] = math.Sin(float64(i))
		x[i] = float64(i)
	}
	xs := ymtn.RSST(x, 5, 4)
	assert.NotNil(xs)

}
