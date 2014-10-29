package ymtn

import (
	"fmt"
	"github.com/hrautila/linalg"
	"github.com/hrautila/linalg/lapack"
	"github.com/hrautila/matrix"
)

func RSST(x []float64, w, n int) []float64 {

	T := len(x)
	stop := w * n
	for t := stop; t < T-stop; t++ {
		//m := n
		//g := 0
		herkel := calcHerkel(x, t, w, n)
		sigma := matrix.FloatZeros(w, n)
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
		//fmt.Printf("Vt: %s\n", vt)
	}
	return nil
}

func calcHerkel(x []float64, t, w, n int) *matrix.FloatMatrix {

	herkel := matrix.FloatZeros(w, n)
	for wi := 0; wi < w; wi++ {
		for ni := 0; ni < n; ni++ {
			i := (t - n + ni) + (1 - w + wi)
			herkel.SetAt(wi, ni, x[i])
		}
	}
	return herkel
}
