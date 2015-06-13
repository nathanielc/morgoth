package stat

import (
	. "github.com/nathanielc/morgoth/stat/fn"
)

func Dirichlet_PDF(α []float64) func(θ []float64) float64 {
	return func(θ []float64) float64 {
		if len(θ) != len(α) {
			return 0
		}
		l := float64(1.0)
		totalα := float64(0)
		for i := 0; i < len(θ); i++ {
			if θ[i] < 0 || θ[i] > 1 {
				return 0
			}
			l *= pow(θ[i], α[i]-1)
			l /= Γ(α[i])
			totalα += α[i]
		}
		l *= Γ(totalα)
		return l
	}
}
func Dirichlet_LnPDF(α []float64) func(x []float64) float64 {
	return func(x []float64) float64 {
		if len(x) != len(α) {
			return negInf
		}
		l := fZero
		totalα := float64(0)
		for i := 0; i < len(x); i++ {
			if x[i] < 0 || x[i] > 1 {
				return negInf
			}
			l += (α[i] - 1) * log(x[i])
			l -= LnΓ(α[i])
			totalα += α[i]
		}
		l += LnΓ(totalα)
		return l
	}
}
func NextDirichlet(α []float64) []float64 {
	x := make([]float64, len(α))
	sum := fZero
	for i := 0; i < len(α); i++ {
		x[i] = NextGamma(α[i], 1.0)
		sum += x[i]
	}
	for i := 0; i < len(α); i++ {
		x[i] /= sum
	}
	return x
}
func Dirichlet(α []float64) func() []float64 {
	return func() []float64 { return NextDirichlet(α) }
}
