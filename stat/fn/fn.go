// Copyright 2012 - 2013 The Fn Authors. All rights reserved. See the LICENSE file.

package fn

// Special math functions.

import (
	"math"
)

const π = float64(math.Pi)
const ln2 = math.Ln2
const lnSqrt2π = 0.918938533204672741780329736406 // log(sqrt(2*pi))
const min64 = math.SmallestNonzeroFloat64         //   DBL_MIN
const eps64 = 1.1102230246251565e-16              // DBL_EPSILON   
const maxExp = 1024.0                             // DBL_MAX_EXP
const sqrt2 = math.Sqrt2
const lnSqrtπd2 = 0.225791352644727432363097614947 // log(sqrt(pi/2))  M_LN_SQRT_PId2

var nan = math.NaN()

var fZero float64 = float64(0.0)
var fOne float64 = float64(1.0)
var iZero int64 = int64(0)
var iOne int64 = int64(1)

var negInf float64 = math.Inf(-1)
var posInf float64 = math.Inf(+1)

// Functions imported from "math"
var abs func(float64) float64 = math.Abs
var floor func(float64) float64 = math.Floor
var ceil func(float64) float64 = math.Ceil
var log func(float64) float64 = math.Log
var log1p func(float64) float64 = math.Log1p
var exp func(float64) float64 = math.Exp
var sqrt func(float64) float64 = math.Sqrt
var pow func(float64, float64) float64 = math.Pow
var atan func(float64) float64 = math.Atan
var tan func(float64) float64 = math.Tan
var sin func(float64) float64 = math.Sin
var trunc func(float64) float64 = math.Trunc
var erf func(float64) float64 = math.Erf
var erfc func(float64) float64 = math.Erfc
var isNaN func(float64) bool = math.IsNaN
var isInf func(float64, int) bool = math.IsInf
var fmod func(float64, float64) float64 = math.Mod

var Γ = math.Gamma
var GammaF = math.Gamma
var sqrt2pi = math.Sqrt(2 * math.Pi)
var logsqrt2pi = math.Log(math.Sqrt(2 * math.Pi))
var LnΓp = LnGammaP
var LnΓpRatio = LnGammaPRatio

var lanczos_coef []float64 = []float64{
	0.99999999999980993,
	676.5203681218851,
	-1259.1392167224028,
	771.32342877765313,
	-176.61502916214059,
	12.507343278686905,
	-0.13857109526572012,
	9.9843695780195716e-6,
	1.5056327351493116e-7}

func isOdd(k float64) bool {
	if k != 2*floor(k/2.0) {
		return true
	}
	return false
}

func isInt(x float64) bool {
	if abs((x)-floor((x)+0.5)) <= 1e-7 {
		return true
	}
	return false
}

// Round to nearest integer
func Round(x float64) float64 {
	var i float64
	f := math.Floor(x)
	c := math.Ceil(x)
	if x-f < c-x {
		i = f
	} else {
		i = c
	}
	return i
}

// Arithmetic mean
func ArithMean(data *Vector) float64 {
	n := data.L
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += data.Get(i)
	}
	return sum / float64(n)
}

// Geometric mean
func GeomMean(data *Vector) float64 {
	n := data.L
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += math.Log(data.Get(i))
	}
	return math.Exp(sum / float64(n))
}

// Harmonic mean
func HarmonicMean(data *Vector) float64 {
	n := data.L
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += 1.0 / data.Get(i)
	}
	return float64(n) / sum
}

// Generalized mean
func GenMean(data *Vector, p float64) float64 {
	n := data.L
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += math.Pow(data.Get(i), p)
	}
	return math.Pow(sum/float64(n), 1/p)
}

// Bernoulli number
// Akiyama–Tanigawa algorithm for Bn
func Bn(n int64) float64 {
	var m int64
	a := make([]float64, n)
	for m = 0; m <= n; m++ {
		a[m] = 1 / float64(m+1)
		for j := m; j >= 1; j-- {
			a[j-1] = float64(j) * (a[j-1] - a[j])
		}
	}
	return a[0] // (which is Bn)
}

// H returns the generalized harmonic number of order n of m. 
func H(n int64, m float64) float64 {
	var i int64
	h := 0.0
	for i = 1; i <= n; i++ {
		h += math.Pow(float64(i), m)
	}
	return h
}

// Generalized harmonic number
func H2(n int64, q, s float64) float64 {
	var i int64
	h := 0.0
	for i = 1; i <= n; i++ {
		h += math.Pow((float64(i) + q), -s)
	}
	return h
}
