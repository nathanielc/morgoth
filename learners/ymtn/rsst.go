package ymtn

import (
	"fmt"
	"github.com/hrautila/linalg"
	"github.com/hrautila/linalg/lapack"
	"github.com/hrautila/matrix"
	//"math"
)

func RSST(x []float64, w, n int) []float64 {

	T := len(x)
	stop := w * n
	for t := stop; t < T-stop; t++ {
		m := n
		g := 0

		u := calcU(x, t, w, n)
		beta, lambda := calcB(x, t, g, w, m)

		fmt.Println("U: ", u)
		fmt.Println("B: ", beta)
		fmt.Println("L: ", lambda)


		

	}
	return nil
}

func calcSHerkel(x []float64, t, w, n int) *matrix.FloatMatrix {

	herkel := matrix.FloatZeros(w, n)
	for wi := 0; wi < w; wi++ {
		for ni := 0; ni < n; ni++ {
			i := (t - n + ni) + (1 - w + wi)
			herkel.SetAt(wi, ni, x[i])
		}
	}
	return herkel
}

func calcl(sigma *matrix.FloatMatrix) int {
	l := sigma.Rows()
	return l
}

func calcGHerkel(x []float64, t, g, w, m int) *matrix.FloatMatrix {
	herkel := matrix.FloatZeros(w, m)
	for wi := 0; wi < w; wi++ {
		for mi := 0; mi < m; mi++ {
			i := (t + g + mi) + (wi)
			herkel.SetAt(wi, mi, x[i])
		}
	}
	return herkel
}


func calcU(x []float64, t, w, n int) *matrix.FloatMatrix {
	herkel := calcSHerkel(x, t, w, n)
	size := w
	if n < w {
		size = n
	}
	sigma := matrix.FloatZeros(size, 1)
	u := matrix.FloatZeros(w, w)
	vt := matrix.FloatZeros(n, n)


	lapack.GesvdFloat(
		herkel,
		sigma,
		u,
		vt,
		&linalg.IOpt{"jobu", linalg.PJobAll},
		&linalg.IOpt{"jobvt", linalg.PJobAll},
	)
	//fmt.Printf("U: %s\n", u)
	//fmt.Printf("S: %s\n", sigma)
	//fmt.Printf("Vt: %s\n", vt)

	l := calcl(sigma)

	sub := matrix.FloatZeros(w, l)
	u.SubMatrix(sub, 0, 0, w, l)

	return sub
}

func calcB(x []float64, t, g, w, m int) (*matrix.FloatMatrix, *matrix.FloatMatrix) {
	herkel := calcGHerkel(x, t, g, w, m)
	size := w
	if m < w {
		size = m
	}
	sigma := matrix.FloatZeros(size, 1)
	u := matrix.FloatZeros(w, w)
	vt := matrix.FloatZeros(m, m)

	lapack.GesvdFloat(
		herkel,
		sigma,
		u,
		vt,
		&linalg.IOpt{"jobu", linalg.PJobAll},
		&linalg.IOpt{"jobvt", linalg.PJobNo},
	)

	l := calcl(sigma)

	sub := matrix.FloatZeros(w, l)
	u.SubMatrix(sub, 0, 0, w, l)

	//Lambdas are the squares of sigmas
	for i := 0; i < size; i++ {
		v := sigma.GetAt(i, 0)
		sigma.SetAt(i, 0, v*v)
	}

	return sub, sigma
}
