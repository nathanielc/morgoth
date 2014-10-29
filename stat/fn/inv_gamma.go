// Incomplete Gamma functions 
// from Cephes Math Library Release 2.8:  June, 2000
// Copyright 1985, 1987, 2000 by Stephen L. Moshier

package fn

import (
	"math"
)

// Complemented incomplete gamma integral
// The function is defined by
//
//
//  IGamC(a,x)   =   1 - IGam(a,x)
//
//                            inf.
//                              -
//                     1       | |  -t  a-1
//               =   -----     |   e   t   dt.
//                    -      | |
//                   | (a)    -
//                             x
//
//
// In this implementation both arguments must be positive.
// The integral is evaluated by either a power series or
// continued fraction expansion, depending on the relative
// values of a and x.
//
// ACCURACY:
//
// Tested at random a, x.
//                a         x                      Relative error:
// arithmetic   domain   domain     # trials      peak         rms
//    IEEE     0.5,100   0,100      200000       1.9e-14     1.7e-15
//    IEEE     0.01,0.5  0,100      200000       1.4e-13     1.6e-15

const (
	MACHEP float64 = 1.11022302462515654042e-16
	MAXLOG float64 = 7.08396418532264106224e2
	big    float64 = 4.503599627370496e15
	biginv float64 = 2.22044604925031308085e-16
)

func IGamC(a, x float64) float64 {
	var ans, ax, c, yc, r, t, y, z, pk, pkm1, pkm2, qk, qkm1, qkm2 float64

	if x <= 0 || a <= 0 {
		return 1.0
	}

	if x < 1.0 || x < a {
		return 1.0 - IGam(a, x)
	}
	ax = a*math.Log(x) - x - LnΓ(a)
	if ax < -MAXLOG {
		panic("IGamC: UNDERFLOW")
		return 0.0
	}
	ax = math.Exp(ax)

	/* continued fraction */
	y = 1.0 - a
	z = x + y + 1.0
	c = 0.0
	pkm2 = 1.0
	qkm2 = x
	pkm1 = x + 1.0
	qkm1 = z * x
	ans = pkm1 / qkm1

	for t > MACHEP {
		c += 1.0
		y += 1.0
		z += 2.0
		yc = y * c
		pk = pkm1*z - pkm2*yc
		qk = qkm1*z - qkm2*yc
		if qk != 0 {
			r = pk / qk
			t = math.Abs((ans - r) / r)
			ans = r
		} else {
			t = 1.0
		}
		pkm2 = pkm1
		pkm1 = pk
		qkm2 = qkm1
		qkm1 = qk
		if math.Abs(pk) > big {
			pkm2 *= biginv
			pkm1 *= biginv
			qkm2 *= biginv
			qkm1 *= biginv
		}
	}
	return (ans * ax)
}

// Incomplete gamma integral
// The function is defined by
//
//                           x
//                            -
//                   1       | |  -t  a-1
//  IGam(a,x)  =   -----     |   e   t   dt.
//                  -      | |
//                 | (a)    -
//                           0
//
//
// In this implementation both arguments must be positive.
// The integral is evaluated by either a power series or
// continued fraction expansion, depending on the relative
// values of a and x.
//
// ACCURACY:
//
//                      Relative error:
// arithmetic   domain     # trials      peak         rms
//    IEEE      0,30       200000       3.6e-14     2.9e-15
//    IEEE      0,100      300000       9.9e-14     1.5e-14
//
// left tail of incomplete gamma function:
//
//          inf.      k
//   a  -x   -       x
//  x  e     >   ----------
//           -     -
//          k=0   | (a+k+1)
//

func IGam(a, x float64) float64 {
	var ans, ax, c, r float64

	if x <= 0 || a <= 0 {
		return 0.0
	}
	if x > 1.0 && x > a {
		return 1.0 - IGamC(a, x)
	}
	// Compute  x**a * exp(-x) / gamma(a)
	ax = a*math.Log(x) - x - LnΓ(a)
	if ax < -MAXLOG {
		panic("IGam: UNDERFLOW")
		return 0.0
	}
	ax = math.Exp(ax)

	// power series
	r = a
	c = 1.0
	ans = 1.0

	for c/ans > MACHEP {
		r += 1.0
		c *= x / r
		ans += c
	}

	return ans * ax / a
}
