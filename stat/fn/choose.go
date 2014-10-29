// Copyright 2012 - 2013 The Fn Authors. All rights reserved. See the LICENSE file.

package fn

// Binomial coefficients.
// FChoose(n, k)   and  LnFChoose(n,k) := log(abs(FChoose(n,k))
// 
// These work for the *generalized* binomial theorem,
// i.e., are also defined for non-integer n  (integer k).
// 
// We use the simple explicit product formula for  k <= k_small_max
// and also have added statements to make sure that the symmetry
// (n \\ k ) == (n \\ n-k)  is preserved for non-negative integer n.

import (
	"math"
)

func lfastchoose(n, k float64) float64 {
	return -log(n+1) - LnB(n-k+1, k+1)
}

func lfastchoose2(n, k float64) (float64, int) {
	// mathematically the same as lfastchoose()
	// less stable typically, but useful if n-k+1 < 0 
	//	r := lgammafn_sign(n-k+1, s_choose)
	r, s_choose := math.Lgamma(n - k + 1)
	p, _ := math.Lgamma(n + 1)
	q, _ := math.Lgamma(k + 1)
	return p - q - r, s_choose
}

// FChoose returns generalized binomial coefficient i.e.,  also defined for non-integer n  (integer k).
func FChoose(n, k float64) float64 {
	const k_small_max = 30.0
	// 30 is somewhat arbitrary: it is on the *safe* side:
	// both speed and precision are clearly improved for k < 30.

	k0 := k
	k = floor(k + 0.5)
	// NaNs propagated correctly
	if isNaN(n) || isNaN(k) {
		return n + k
	}
	if abs(k-k0) > 1e-7 {
		panic("k must be integer valued float64")
	}
	if k < k_small_max {
		if n-k < k && n >= 0 && isInt(n) { //  symmetry
			k = n - k
		}
		if k < 0 {
			return 0
		}
		if k == 0 {
			return 1
		}

		// else: k >= 1
		r := n

		for j := 2.0; j <= k; j++ {
			r *= (n - j + 1) / j
		}

		if isInt(n) {
			return floor(r + 0.5)
		} else {
			return floor(r)
		}
	}
	// else: k >= k_small_max
	if n < 0 {
		r := FChoose(-n+k-1, k)
		if isOdd(k) {
			r = -r
		}
		return r
	} else if isInt(n) {
		n = floor(n + 0.5)
		if n < k {
			return 0
		}
		if n-k < k_small_max { //  symmetry
			return FChoose(n, n-k)
		}

		return floor(exp(lfastchoose(n, k)) + 0.5)
	}

	// else non-integer n >= 0
	if n < k-1 {
		r, s_choose := lfastchoose2(n, k)
		return float64(s_choose) * exp(r)
	}
	return exp(lfastchoose(n, k))
}

// LnFChoose returns the natural logarithm of the generalized binomial coefficient i.e.,  also defined for non-integer n  (integer k).
func LnFChoose(n, k float64) float64 {
	k0 := k
	k = floor(k + 0.5)
	// NaNs propagated correctly
	if isNaN(n) || isNaN(k) {
		return n + k
	}
	if abs(k-k0) > 1e-7 {
		panic("k must be integer valued float64")
	}
	if k < 2 {
		if k < 0 {
			return negInf
		}
		if k == 0 {
			return 0
		}
		/* else: k == 1 */
		return log(abs(n))
	}
	/* else: k >= 2 */
	if n < 0 {
		return LnFChoose(-n+k-1, k)
	} else if isInt(n) {

		n = floor(n + 0.5)
		if n < k {
			return negInf
		}
		/* k <= n */
		if n-k < 2 {
			return LnFChoose(n, n-k) /* <- Symmetry */
		}
		/* else: n >= k+2 */
		return lfastchoose(n, k)
	}
	/* else non-integer n >= 0 : */
	if n < k-1 {
		r, _ := lfastchoose2(n, k)
		return r
	}
	return lfastchoose(n, k)
}

// Integer versions.

// Choose returns binomial coefficient for integer n  and k.
func Choose(n int64, i int64) int64 {
	smaller := i
	if n-i < smaller {
		smaller = n - i
	}
	return PartialFactInt(n, smaller) / FactInt(smaller)
}

// LnChoose returns the natural logarithm of the binomial coefficient for integer n  and k.
func LnChoose(n int64, i int64) float64 {
nn:= float64(n)
ii:= float64(i)
	smaller := ii
	if nn-ii < smaller {
		smaller = nn - ii
	}
	return LnPartialFact(nn, smaller) - LnFact(smaller)
}
