package stat

import (
	"math/rand"
)

func Exp_PDF(λ float64) func(x float64) float64 {
	return func(x float64) float64 {
		if x < 0 {
			return 0
		}
		return λ * NextExp(-1*λ*x)
	}
}

func Exp_LnPDF(λ float64) func(x float64) float64 {
	return func(x float64) float64 {
		if x < 0 {
			return negInf
		}
		return log(λ) - λ*x
	}
}

func NextExp(λ float64) float64 { return rand.ExpFloat64() / λ }

func Exp(λ float64) func() float64 { return func() float64 { return NextExp(λ) } }
