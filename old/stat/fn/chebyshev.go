// Copyright 2012 - 2013 The Fn Authors. All rights reserved. See the LICENSE file.

package fn

//The Chebyshev series.

func chebyshev_init(dos []float64, nos int, eta float64) int {
	if nos < 1 {
		return 0
	}
	err := 0.0
	i := 0 /* just to avoid compiler warnings */
	for ii := 1; ii <= nos; ii++ {
		i = nos - ii
		err += abs(dos[i])
		if err > eta {
			return i
		}
	}
	return i
}

func chebyshev_eval(x float64, a []float64, n int) float64 {

	if n < 1 || n > 1000 {
		return nan
	}

	if x < -1.1 || x > 1.1 {
		return nan
	}
	twox := x * 2
	b1 := 0.0
	b2 := b1
	b0 := 0.0
	for i := 1; i <= n; i++ {
		b2 = b1
		b1 = b0
		b0 = twox*b1 - b2 + a[n-i]
	}
	return (b0 - b2) * 0.5
}
