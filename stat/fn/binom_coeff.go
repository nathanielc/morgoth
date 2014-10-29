// Copyright 2012 - 2013 The Fn Authors. All rights reserved. See the LICENSE file.

package fn

//The Binomial coefficient.

// Binomial coefficient (in combinatorics, it gives the number of ways, disregarding order, 
// that k objects can be chosen from among n objects; more formally, the number of 
// k-element subsets (or k-combinations) of an n-element set)

// Same as Choose(n, k), LnChoose(n, k)

import (
	"math"
)

func BinomCoeff(n, k int64) float64 {
	if k == 0 {
		return 1
	}
	if n == 0 {
		return 0
	}
	// if n, k are small, use recursive formula
	if n < 10 && k < 10 {
		return BinomCoeff(n-1, k-1) + BinomCoeff(n-1, k)
	}

	// else, use factorial formula
	//	fmt.Println(LnFactBig(n), LnFactBig(k), LnFactBig(n-k))
	return Round(math.Exp(LnFactBig(float64(n)) - LnFactBig(float64(k)) - LnFactBig(float64(n-k))))
}

func LnBinomCoeff(n, k float64) float64 {
	if k == 0 {
		return math.Log(1)
	}
	if n == 0 {
		panic("n == 0")
	}
	if n < 10 && k < 10 {
nn:= int64(n)
kk:= int64(k)

		return math.Log(BinomCoeff(nn, kk))
	}

	// else, use factorial formula
	return LnFactBig(n) - LnFactBig(k) - LnFactBig(n-k)
}
