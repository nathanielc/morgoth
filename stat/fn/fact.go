// Copyright 2012 - 2013 The Fn Authors. All rights reserved. See the LICENSE file.

package fn

// Factorials.

import (
	"math"
)

func Fact(n int64) float64 {
	var i int64
	if n < 0 {
		return nan
	}
	if n < 101 {
		return factorial[n]
	} //else 

	f := factorial[100]
	for i = 101; i <= n; i++ {
		f *= float64(i)
	}
	return f
}

//FactInt(n) = n*FactInt(n-1)
func FactInt(n int64) int64 {
	return PartialFactInt(n, 0)
}


func LnFact(nn float64) float64 {

n:= trunc(nn)
	var i float64

	switch {
	case n < 0:
		return nan
	case n < 101:
		return log(factorial[int(n)])
	case n < 10000000:
		f := log(factorial[100])
		for i = 101; i <= n; i++ {
			f += log(i)
		}
		return f
	default: // use gamma approximation
		return LnFactBig(n)
	}
	return nan // should never happen
}

//LnFactBig(n) = Gamma(n+1)
func LnFactBig(n float64) float64 {
	n= trunc(n)
	return LnΓ(n + 1)
}

//PartialFactInt returns Fact(n)/Fact(m)
func PartialFactInt(n int64, m int64) int64 {
	if n == m {
		return 1
	}
	return n * PartialFactInt(n-1, m)
}

//LnPartialFact returns LnFact(n)-LnFact(m)
func LnPartialFact(n, m float64) float64 {
	if n == m {
		return 0
	}
	return math.Log(n) + LnPartialFact(n-1, m)
}
