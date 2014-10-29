package stat

import (
	"math/rand"
)

func Range_PMF(n int64) func(i int64) float64 {
	return func(i int64) float64 {
		return fOne / float64(n)
	}
}
func LnRange_PMF(n int64) func(i int64) float64 {
	return func(i int64) float64 {
		return -log(float64(n))
	}
}
func NextRange(n int64) int64 {
	return rand.Int63n(n)
}
func Range(n int64) func() int64 {
	return func() int64 {
		return NextRange(n)
	}
}
