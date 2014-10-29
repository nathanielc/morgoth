// F-distribution, alias Fisher-Snedecor distribution
package stat

import (
	. "github.com/nvcook42/morgoth/stat/fn"
	"fmt"
)

func F_PDF(d1 float64, d2 float64) func(x float64) float64 {
	normalization := 1 / B(d1/2, d2/2)
	return func(x float64) float64 {
		return normalization * sqrt(pow(d1*x, d1)*pow(d2, d2)/pow(d1*x+d2, d1+d2)) / x
	}
}
func F_LnPDF(d1 float64, d2 float64) func(x float64) float64 {
	normalization := -LnB(d1/2, d2/2)
	return func(x float64) float64 {
		return normalization + log(d1*x)*d1/2 + log(d2)*d2/2 - log(d1*x+d2)*(d1+d2)/2 - log(x)
	}
}
func NextF(d1 int64, d2 int64) float64 {
	return (NextXsquare(d1) * float64(d2)) / (NextXsquare(d2) * float64(d1))
}
func F(d1 int64, d2 int64) func() float64 {
	return func() float64 {
		return NextF(d1, d2)
	}
}

// CDF of F-distribution
func F_CDF(df1, df2 float64) func(x float64) float64 {
	return func(x float64) float64 {
		y := df1 * x / (df1*x + df2)
		return BetaIncReg(df1/2.0, df2/2.0, y)
	}
}

// Value of CDF of F-distribution at x
func F_CDF_At(df1, df2, x float64) float64 {
	cdf := F_CDF(df1, df2)
	return cdf(x)
}

// Inverse CDF (Quantile) function of F-distribution
func F_InvCDF(df1, df2 float64) func(p float64) float64 {
	return func(p float64) float64 {
		if p < 0.0 {
			panic(fmt.Sprintf("p < 0"))
		}
		if p > 1.0 {
			panic(fmt.Sprintf("p > 1.0"))
		}
		if df1 < 1.0 {
			panic(fmt.Sprintf("df1 < 1"))
		}
		if df2 < 1.0 {
			panic(fmt.Sprintf("df2 < 1"))
		}

		return ((1/BetaInv_CDF_For(df2/2, df1/2, 1-p) - 1) * df2 / df1)
	}
}

// Value of the inverse CDF of F-distribution for probability p
func F_InvCDF_For(df1, df2, p float64) float64 {
	cdf := F_InvCDF(df1, df2)
	return cdf(p)
}
