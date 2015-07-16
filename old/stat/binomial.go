package stat

import (
	. "github.com/nathanielc/morgoth/stat/fn"
)

// Probability Mass Function for the Binomial distribution
func Binomial_PMF(ρ float64, n int64) func(i int64) float64 {
	return func(i int64) float64 {
		p := pow(ρ, float64(i)) * pow(1-ρ, float64(n-i))
		p *= Γ(float64(n+1)) / (Γ(float64(i+1)) * Γ(float64(n-i+1)))
		return p
	}
}

func Binomial_PMF_At(ρ float64, n, k int64) float64 {
	pmf := Binomial_PMF(ρ, n)
	return pmf(k)
}

// Natural logarithm of Probability Mass Function for the Binomial distribution
func Binomial_LnPMF(ρ float64, n int64) func(i int64) float64 {
	return func(i int64) float64 {
		p := log(ρ)*float64(i) + log(1-ρ)*float64(n-i)
		p += LnΓ(float64(n+1)) - LnΓ(float64(i+1)) - LnΓ(float64(n-i+1))
		return p
	}
}

func NextBinomial(ρ float64, n int64) (result int64) {
	for i := int64(0); i <= n; i++ {
		result += NextBernoulli(ρ)
	}
	return
}

func Binomial(ρ float64, n int64) func() int64 {
	return func() int64 { return NextBinomial(ρ, n) }
}

// Cumulative Distribution Function for the Binomial distribution, trivial implementation
func Binomial_CDF_trivial(ρ float64, n int64) func(k int64) float64 {
	return func(k int64) float64 {
		var p float64 = 0
		var i int64
		pmf := Binomial_PMF(ρ, n)
		for i = 0; i <= k; i++ {
			p += pmf(i)
		}
		return p
	}
}

// Cumulative Distribution Function for the Binomial distribution
func Binomial_CDF(ρ float64, n int64) func(k int64) float64 {
	return func(k int64) float64 {
		p := Beta_CDF_At((float64)(n-k), (float64)(k+1), 1-ρ)
		return p
	}
}

func Binomial_CDF_At(ρ float64, n, k int64) float64 {
	cdf := Binomial_CDF(ρ, n)
	return cdf(k)
}
