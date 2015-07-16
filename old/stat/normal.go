package stat

import (
	"math"
	"math/rand"
)

func rat_eval(a []float64, na int64, b []float64, nb int64, x float64) float64 {
	var (
		i, j    int64
		u, v, r float64
	)

	u = a[na-1]

	for i = na - 1; i > 0; i-- {
		u = x*u + a[i-1]
	}

	v = b[nb-1]

	for j = nb - 1; j > 0; j-- {
		v = x*v + b[j-1]
	}

	r = u / v

	return r
}

func small(q float64) float64 {
	var (
		a = []float64{3.387132872796366608, 133.14166789178437745, 1971.5909503065514427, 13731.693765509461125, 45921.953931549871457, 67265.770927008700853, 33430.575583588128105, 2509.0809287301226727}

		b = []float64{1.0, 42.313330701600911252, 687.1870074920579083, 5394.1960214247511077, 21213.794301586595867, 39307.89580009271061,
			28729.085735721942674, 5226.495278852854561}

		r, x float64
	)
	r = 0.180625 - q*q
	x = q * rat_eval(a, 8, b, 8, r)

	return x
}

func intermediate(r float64) float64 {
	var (
		a = []float64{1.42343711074968357734, 4.6303378461565452959, 5.7694972214606914055, 3.64784832476320460504,
			1.27045825245236838258, 0.24178072517745061177, 0.0227238449892691845833, 7.7454501427834140764e-4}

		b = []float64{1.0, 2.05319162663775882187, 1.6763848301838038494, 0.68976733498510000455, 0.14810397642748007459, 0.0151986665636164571966,
			5.475938084995344946e-4, 1.05075007164441684324e-9}

		x float64
	)
	x = rat_eval(a, 8, b, 8, (r - 1.6))
	return x
}

func tail(r float64) float64 {

	var (
		a = []float64{6.6579046435011037772, 5.4637849111641143699, 1.7848265399172913358, 0.29656057182850489123,
			0.026532189526576123093, 0.0012426609473880784386, 2.71155556874348757815e-5, 2.01033439929228813265e-7}

		b = []float64{1.0, 0.59983220655588793769, 0.13692988092273580531, 0.0148753612908506148525, 7.868691311456132591e-4, 1.8463183175100546818e-5, 1.4215117583164458887e-7, 2.04426310338993978564e-15}

		x float64
	)
	x = rat_eval(a, 8, b, 8, (r - 5.0))
	return x
}

func Normal_PDF(μ float64, σ float64) func(x float64) float64 {
	normal_normalizer := 0.3989422804014327 / σ
	return func(x float64) float64 { return normal_normalizer * exp(-1*(x-μ)*(x-μ)/(2*σ*σ)) }
}

func Normal_LnPDF(μ float64, σ float64) func(x float64) float64 {
	ln_normal_normalizer := -0.91893853320467267 - log(σ)
	return func(x float64) float64 { return ln_normal_normalizer - (x-μ)*(x-μ)/(2*σ*σ) }
}

func NextNormal(μ float64, σ float64) float64 { return rand.NormFloat64()*σ + μ }

func Normal(μ, σ float64) func() float64 {
	return func() float64 { return NextNormal(μ, σ) }
}

// Cumulative Distribution Function for the Normal distribution
func Normal_CDF(μ, σ float64) func(x float64) float64 {
	return func(x float64) float64 { return ((1.0 / 2.0) * (1 + math.Erf((x-μ)/(σ*math.Sqrt2)))) }
}

// Inverse CDF of Normal distribution for probability p //// ??? only *sigma? from GSL
func NormalInv_CDF_For(p, sigma float64) float64 {
	return sigma * Z_InvCDF_For(p)
}

// Probability Density Function for the Standard Normal distribution
func Z_PDF() func(float64) float64 {
	return Normal_PDF(0, 1)
}

// Cumulative Distribution Function for the Standard Normal distribution
func Z_CDF() func(float64) float64 {
	return Normal_CDF(0, 1)
}

// Probability Density of the Standard Normal distribution at x
func Z_PDF_At(x float64) float64 {
	pdf := Normal_PDF(0, 1)
	return pdf(x)
}

// Cumulative Probability of the Standard Normal distribution at x
func Z_CDF_At(x float64) float64 {
	cdf := Normal_CDF(0, 1)
	return cdf(x)
}

// Inverse CDF of Standard Normal distribution for probability p
func Z_InvCDF_For(p float64) float64 {

	var r, x, pp, dp float64

	dp = p - 0.5
	switch {
	case p == 1.0:
		return math.MaxFloat64
	case p == 0.0:
		return -math.MaxFloat64
	}
	if math.Abs(dp) <= 0.425 {
		x = small(dp)
		return x
	}
	if p < 0.5 {
		pp = p
	} else {
		pp = 1.0 - p
	}
	r = math.Sqrt(-math.Log(pp))
	if r <= 5.0 {
		x = intermediate(r)
	} else {
		x = tail(r)
	}
	if p < 0.5 {
		return -x
	}
	return x
}
