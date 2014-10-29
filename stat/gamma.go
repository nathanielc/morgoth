// Gamma distribution
// k > 0		shape parameter
// θ (Theta) > 0	scale parameter

package stat

import (
	. "github.com/nvcook42/morgoth/stat/fn"
	"fmt"
	"math"
)

/* did not pass test, so commented out
// Probability density function
func Gamma_PDF(α float64, λ float64) func(x float64) float64 {
	expPart := Exp_PDF(λ)
	return func(x float64) float64 {
		if x < 0 {
			return 0
		}
		return expPart(x) * pow(λ*x, α-1) / Γ(α)
	}
}
*/

// Probability density function
func Gamma_PDF(k float64, θ float64) func(x float64) float64 {
	return func(x float64) float64 {
		if x < 0 {
			return 0
		}
		return pow(x, k-1) * exp(-x/θ) / (Γ(k) * pow(θ, k))
	}
}

// Natural logarithm of the probability density function
func Gamma_LnPDF(α float64, λ float64) func(x float64) float64 {
	expPart := Exp_LnPDF(λ)
	return func(x float64) float64 {
		if x < 0 {
			return negInf
		}
		return expPart(x) + (α-1)*log(λ*x) - LnΓ(α)
	}
}

// Random value drawn from the distribution
func NextGamma(α float64, λ float64) float64 {
	//if α is a small integer, this way is faster on my laptop
	if α == float64(int64(α)) && α <= 15 {
		x := NextExp(λ)
		for i := 1; i < int(α); i++ {
			x += NextExp(λ)
		}
		return x
	}

	if α < 0.75 {
		return RejectionSample(Gamma_PDF(α, λ), Exp_PDF(λ), Exp(λ), 1)
	}

	//Tadikamalla ACM '73
	a := α - 1
	b := 0.5 + 0.5*sqrt(4*α-3)
	c := a * (1 + b) / b
	d := (b - 1) / (a * b)
	s := a / b
	p := 1.0 / (2 - exp(-s))
	var x, y float64
	for i := 1; ; i++ {
		u := NextUniform()
		if u > p {
			var e float64
			for e = -log((1 - u) / (1 - p)); e > s; e = e - a/b {
			}
			x = a - b*e
			y = a - x
		} else {
			x = a - b*log(u/p)
			y = x - a
		}
		u2 := NextUniform()
		if log(u2) <= a*log(d*x)-x+y/b+c {
			break
		}
	}
	return x / λ
}

func Gamma(α float64, λ float64) func() float64 {
	return func() float64 { return NextGamma(α, λ) }
}

// Cumulative distribution function, analytic solution, did not pass some tests!
func Gamma_CDF(k float64, θ float64) func(x float64) float64 {
	return func(x float64) float64 {
		if k < 0 || θ < 0 {
			panic(fmt.Sprintf("k < 0 || θ < 0"))
		}
		if x < 0 {
			return 0
		}
		return Iγ(k, x/θ) / Γ(k)
	}
}

// Cumulative distribution function, for integer k only
func Gamma_CDFint(k int64, θ float64) func(x float64) float64 {
	return func(x float64) float64 {
		if k < 0 || θ < 0 {
			panic(fmt.Sprintf("k < 0 || θ < 0"))
		}
		if x < 0 {
			return 0
		}
		return Iγint(k, x/θ) / Γ(float64(k))
	}
}

/*
// Cumulative distribution function, using gamma incomplete integral  DOES NOT WORK !!!
func Gamma_CDF(k float64, θ float64) func(x float64) float64 {
	return func(x float64) float64 {
		if k < 0 || θ < 0 {
			panic(fmt.Sprintf("k < 0 || θ < 0"))
		}
		if x < 0 {
			return 0
		}
		return IGam(θ, k*x)
	}
}
*/

// Value of the probability density function at x
func Gamma_PDF_At(k, θ, x float64) float64 {
	pdf := Gamma_PDF(k, θ)
	return pdf(x)
}

// Value of the cumulative distribution function at x
func Gamma_CDF_At(k, θ, x float64) float64 {
	cdf := Gamma_CDF(k, θ)
	return cdf(x)
}

// Inverse CDF (Quantile) function
func Gamma_InvCDF(k float64, θ float64) func(x float64) float64 {
	return func(x float64) float64 {
		var eps, y_new, h float64
		eps = 1e-4
		y := k * θ
		y_old := y
	L:
		for i := 0; i < 100; i++ {
			h = (Gamma_CDF_At(k, θ, y_old) - x) / Gamma_PDF_At(k, θ, y_old)
			y_new = y_old - h
			if y_new <= eps {
				y_new = y_old / 10
				h = y_old - y_new
			}
			if math.Abs(h) < eps {
				break L
			}
			y_old = y_new
		}
		return y_new
	}
}

// Value of the inverse CDF for probability p
func Gamma_InvCDF_For(k, θ, p float64) float64 {
	cdf := Gamma_InvCDF(k, θ)
	return cdf(p)
}
