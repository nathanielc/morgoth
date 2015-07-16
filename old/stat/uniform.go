// Uniform (Flat) distribution

package stat

import (
	"math/rand"
)

func Uniform_PDF() func(x float64) float64 {
	return func(x float64) float64 {
		if 0 <= x && x <= 1 {
			return 1
		}
		return 0
	}
}

func Uniform_LnPDF() func(x float64) float64 {
	return func(x float64) float64 {
		if 0 <= x && x <= 1 {
			return 0
		}
		return negInf
	}
}

var NextUniform func() float64 = rand.Float64

func Uniform() func() float64 { return NextUniform }
