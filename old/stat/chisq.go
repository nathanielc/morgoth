// Chi-Squared distribution

package stat

import (
	. "github.com/nathanielc/morgoth/stat/fn"
)

func Xsquare_PDF(n int64) func(x float64) float64 {
	k := float64(n) / 2
	normalization := pow(0.5, k) / Γ(k)
	return func(x float64) float64 {
		return normalization * pow(x, k-1) * NextExp(-x/2)
	}
}

func Xsquare_LnPDF(n int64) func(x float64) float64 {
	k := float64(n) / 2
	normalization := log(0.5)*k - LnΓ(k)
	return func(x float64) float64 {
		return normalization + log(x)*(k-1) - x/2
	}
}

//Xsquare(n) => sum of n N(0,1)^2
func NextXsquare(n int64) (x float64) {
	for i := iZero; i < n; i++ {
		n := NextNormal(0, 1)
		x += n * n
	}
	return
}

func Xsquare(n int64) func() float64 {
	return func() float64 {
		return NextXsquare(n)
	}
}

//Cumulative density function of the Chi-Squared distribution
func Xsquare_CDF(n int64) func(p float64) float64 {
	return func(p float64) float64 {
		return Γr(float64(n)/2, p/2)
	}
}

//Inverse CDF (Quantile) function of the Chi-Squared distribution
func Xsquare_InvCDF(n int64) func(p float64) float64 {
	return func(p float64) float64 {
		//return Gamma_InvCDF_At(n/2, 2, p)  to be implemented
		return Gamma_InvCDF_For(float64(n)/2, 2, p)
	}
}
