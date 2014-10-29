package stat

import (
	. "github.com/nvcook42/morgoth/stat/fn"
)

func Multinomial_PMF(θ []float64, n int64) func(x []int64) float64 {
	return func(x []int64) float64 {
		if len(x) != len(θ) {
			return 0
		}
		l := fOne
		totalx := iZero
		for i := 0; i < len(x); i++ {
			l *= pow(θ[i], float64(x[i]))
			l /= Γ(float64(x[i] + 1))
			totalx += x[i]
		}
		if totalx != n {
			return 0
		}
		l *= Γ(float64(totalx + 1))
		return l
	}
}
func Multinomial_LnPMF(θ []float64, n int64) func(x []int64) float64 {
	return func(x []int64) float64 {
		if len(x) != len(θ) {
			return negInf
		}
		l := fZero
		totalx := iZero
		for i := 0; i < len(x); i++ {
			l += log(θ[i]) * float64(x[i])
			l -= LnΓ(float64(x[i] + 1))
			totalx += x[i]
		}
		if totalx != n {
			return negInf
		}
		l += LnΓ(float64(totalx + 1))
		return l
	}
}
func NextMultinomial(θ []float64, n int64) []int64 {
	x := make([]int64, len(θ))
	chooser := Choice(θ)
	for i := iZero; i < n; i++ {
		x[chooser()]++
	}
	return x
}
func Multinomial(θ []float64, n int64) func() []int64 {
	return func() []int64 {
		return NextMultinomial(θ, n)
	}
}
