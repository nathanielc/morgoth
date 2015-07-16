package stat

import (
	"fmt"
	. "github.com/nathanielc/morgoth/stat/fn"
	"math"
	"math/rand"
	"testing"
	"time"
)

var Seed func(int64) = rand.Seed

func XTestDir(t *testing.T) {
	α := []float64{4, 5, 6}
	dgen := Dirichlet(α)
	counts := [3]int{0, 0, 0}
	const total = 150000
	for i := 0; i < total; i++ {
		θ := dgen()
		v := NextChoice(θ)
		counts[v]++
	}
	fmt.Printf("%v\n", counts)
}

func TestNullWeights(t *testing.T) {
	n := int64(10)
	weights := make([]float64, n)
	m := NextChoice(weights)
	if n != m {
		t.Error()
	}
}

func TestLnGamma(t *testing.T) {
	acc := 0.0000001
	check := func(x, y float64) bool {
		if false {
			return x == y
		}
		return math.Abs(x-y) < acc
	}
	for i := 0; i < 100; i++ {
		x := NextGamma(10, 10)
		g1 := LnΓ(x)
		g2, _ := math.Lgamma(x)
		if !check(g1, g2) {
			t.Error(fmt.Sprintf("For %v: %v vs %v", x, g1, g2))
		}
	}
	//var start int64
	Seed(10)
	start := time.Now()
	for i := 0; i < 1e6; i++ {
		x := NextGamma(10, 10)
		math.Lgamma(x)
	}
	now := time.Now()
	duration2 := float64(now.Sub(start)) / 1e9

	//duration2 := float64(time.Now()-start) / 1e9
	Seed(10)
	start = time.Now()
	for i := 0; i < 1e6; i++ {
		x := NextGamma(10, 10)
		LnΓ(x)
	}
	now = time.Now()
	duration1 := float64(now.Sub(start)) / 1e9
	fmt.Printf("Mine was %f\nTheirs was %f\n", duration1, duration2)
}

func XTestGen(t *testing.T) {
	fmt.Printf("NextUniform => %f\n", NextUniform())
	fmt.Printf("NextExp => %f\n", NextExp(1.5))
	fmt.Printf("NextGamma => %f\n", NextGamma(.3, 1))
	fmt.Printf("NextNormal => %f\n", NextNormal(0, 1))
	fmt.Printf("NextRange => %d\n", NextRange(10))
	fmt.Printf("NextChoice => %d\n", NextChoice([]float64{.3, .3, .4}))
	fmt.Printf("NextMultinomial => %v\n",
		NextMultinomial([]float64{.3, .3, .4}, 100))
	fmt.Printf("NextDirichlet => %v\n", NextDirichlet([]float64{.3, .3, .4}))
	fmt.Printf("NextBernoulli => %d\n", NextBernoulli(.5))
	fmt.Printf("NextGeometric => %d\n", NextGeometric(.5))
	fmt.Printf("NextBinomial => %d\n", NextBinomial(.5, 10))
	fmt.Printf("NextXsquare => %f\n", NextXsquare(3))
	fmt.Printf("NextNegativeBinomial => %d\n", NextNegativeBinomial(.5, 10))
	fmt.Printf("NextStudentsT => %f\n", NextStudentsT(7))
	fmt.Printf("NextF => %f\n", NextF(7, 3))
	/*	fmt.Printf("NextWishart => %v\n",

			NextWishart(100, matrix.MakeDenseMatrixStacked([][]float64{[]float64{1, 0}, []float64{0, 1}})))
		fmt.Printf("NextInverseWishart => %v\n",
			NextInverseWishart(100, matrix.MakeDenseMatrixStacked([][]float64{[]float64{1, 0}, []float64{0, 1}})))
	*/
}

/*
// test for Binomial p confidence interval  // failed due to some unknown bug
func TestBinomP_CI(t *testing.T) {
	fmt.Println("test for Binomial p confidence interval")
	var n int64
	var k, nn, p, alpha, low, high, low2, high2 float64
	n=30
	nn=float64(n)
	p=0.1
	k=nn*p

	alpha=0.1
	low2=0.04
	high2=0.21
	low, high =  Binom_p_ConfI(n, p, alpha)
	fmt.Println(low, " = ", low2, "\t", high, " = ",  high2)
}
*/
