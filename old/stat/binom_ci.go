package stat

// One-sided (frequentist) Confidence Intervals for Observed "Nonconforming" Units in a Random Sample
// Source: Hahn, G. J., and W. Q. Meeker, "Statistical Intervals / A Guide for Practitioners," J. Wiley & Sons, New York.  1991.

func Binom_p_ConfI(n int64, p, alpha float64) (float64, float64) {

	/*
		Alpha	100(1-alpha) is the confidence
		n	Sample size
		p	Observed proportion
		lCL	Lower confidence limit
		uCL	Upper confidence limit
	*/

	var nn, k, lCL, uCL float64
	nn = float64(n)
	k = nn * p
	if k <= 0 {
		lCL = 0.0
	} else {
		lCL = 1.0 / (1.0 + (nn-k+1)*F_InvCDF_For(alpha, 2*nn-2*k+2, 2*k)/k)
	}

	if k >= nn {
		uCL = 1.0
	} else {
		uCL = 1.0 / (1.0 + (nn-k)/((k+1)*F_InvCDF_For(alpha, 2*k+2, 2*nn-2*k)))
	}
	return lCL, uCL
}
