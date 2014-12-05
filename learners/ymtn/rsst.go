package ymtn

import (
	"github.com/golang/glog"
	"github.com/nvcook42/linalg"
	"github.com/nvcook42/linalg/blas"
	"github.com/nvcook42/linalg/lapack"
	"github.com/nvcook42/matrix"
	"math"
)

const (
	small = 1e-9
)

func RSST(x []float64, w, n int) []float64 {

	glog.V(2).Infoln("x: ", x)

	T := len(x)
	m := n
	g := 0

	//Calc change scores of x
	changeScores := make([]float64, T)
	for t := w + n; t < T-w-m-g-1; t++ {

		past := calcPast(x, t, w, n)
		future, eigenValues := calcFuture(x, t, g, w, m)

		glog.V(2).Infoln("Past: ", past)
		glog.V(2).Infoln("Future: ", future)
		glog.V(2).Infoln("Eigen Values: ", eigenValues)

		//Only calc changescores if our eigenvalues are not small
		if eigenValues.GetAt(0, 0) > small {
			changeScores[t] = calcChangeScore(past, future)
		} // else changescores[t] = 0
	}
	glog.V(1).Infoln("Change Scores: ", changeScores)

	//Weight each score by it's past and future
	width := int(float64(w)/2.0 + 0.5)
	start := w + n
	if start < width {
		start = width
	}
	stop := T - width - 1
	glog.V(3).Infoln("width: ", width)
	glog.V(3).Infoln("start: ", start)
	glog.V(3).Infoln("stop: ", stop)

	weighted := make([]float64, T)
	max := 0.0
	for t := start; t < stop; t++ {
		pastMean, pastVar := calcMeanVar(changeScores[t-width : t])
		futureMean, futureVar := calcMeanVar(changeScores[t : t+width])
		glog.V(4).Infoln("t, pm, pv, fm, fv", t, pastMean, pastVar, futureMean, futureVar)
		score :=
			changeScores[t] *
				math.Abs(pastMean-futureMean) *
				math.Abs(math.Sqrt(pastVar)-math.Sqrt(futureVar))
		weighted[t] = score
		if score > max {
			max = score
		}
	}
	glog.V(1).Infoln("Weighted: ", weighted)

	if max == 0 {
		//Trivial zero case, we can return early
		return weighted
	}

	//Keep only local trimmed
	trimmed := make([]float64, T)
	for i, v := range weighted {
		if (i == 0 || v > weighted[i-1]) &&
			(i == T-1 || v > weighted[i+1]) {
			//Keep and normalize
			trimmed[i] = v / max
		}
	}
	glog.V(1).Infoln("Trimmed: ", trimmed)

	return trimmed
}

func calcMeanVar(x []float64) (float64, float64) {
	l := float64(len(x))
	sum := 0.0
	for _, v := range x {
		sum += v
	}
	mean := sum / l
	varSum := 0.0
	for _, v := range x {
		diff := v - mean
		varSum += diff * diff
	}

	variance := varSum / l

	glog.V(4).Infoln(l, sum, mean, variance, x)
	return mean, variance
}

func calcPastHerkel(x []float64, t, w, n int) *matrix.FloatMatrix {

	herkel := matrix.FloatZeros(w, n)
	for wi := 0; wi < w; wi++ {
		for ni := 0; ni < n; ni++ {
			i := (t - n + ni) + (1 - w + wi)
			herkel.SetAt(wi, ni, x[i])
		}
	}
	return herkel
}

//Find the number of eigenvectors to use based on
// the corner of the accumlative sum of sigma
func calcNumEigenValues(sigma *matrix.FloatMatrix) int {
	glog.V(3).Infof("sigma: %v", sigma)
	l := sigma.Rows()
	if sigma.GetAt(0, 0) < small || l == 1 {
		return 1
	}

	sum := 0.0
	cumsum := make([]float64, l)
	for i := 0; i < l; i++ {
		sum += sigma.GetAt(i, 0)
		cumsum[i] = sum
	}

	glog.V(3).Infof("sum: %f", sum)
	glog.V(3).Infof("cumsum: %v", cumsum)

	i := 0
	for ; cumsum[i]/sum < 0.95; i++ {
	}
	//Plus 1 since indexes are zero-based and we want count
	i++

	glog.V(2).Infof("NumEV: %d", i)
	return i
}

func calcFutureHerkel(x []float64, t, g, w, m int) *matrix.FloatMatrix {
	herkel := matrix.FloatZeros(w, m)
	for wi := 0; wi < w; wi++ {
		for mi := 0; mi < m; mi++ {
			i := (t + g + mi) + (wi)
			herkel.SetAt(wi, mi, x[i])
		}
	}
	return herkel
}

func calcPast(x []float64, t, w, n int) *matrix.FloatMatrix {
	herkel := calcPastHerkel(x, t, w, n)
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

	l := calcNumEigenValues(sigma)
	glog.V(3).Infoln("lp: ", l)

	sub := matrix.FloatZeros(w, l)
	u.SubMatrix(sub, 0, 0, w, l)

	return sub
}

func calcFuture(x []float64, t, g, w, m int) (*matrix.FloatMatrix, *matrix.FloatMatrix) {
	herkel := calcFutureHerkel(x, t, g, w, m)
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

	//Eigen values are the squares of sigmas
	eigenValues := matrix.FloatZeros(sigma.Rows(), 1)
	for i := 0; i < sigma.Rows(); i++ {
		v := sigma.GetAt(i, 0)
		eigenValues.SetAt(i, 0, v*v)
	}

	//For future, use eigenvalues instead of sigma
	l := calcNumEigenValues(eigenValues)
	glog.V(3).Infoln("lf: ", l)

	sub := matrix.FloatZeros(w, l)
	u.SubMatrix(sub, 0, 0, w, l)

	eigenValues.SubMatrix(eigenValues, 0, 0, l, 1)

	return sub, eigenValues
}

func calcChangeScore(past, future *matrix.FloatMatrix) float64 {

	w := past.Rows()
	lf := future.Cols()
	b := matrix.FloatZeros(w, 1)
	v := matrix.FloatZeros(past.Cols(), 1)
	csSum := 0.0
	for i := 0; i < lf; i++ {
		future.SubMatrix(b, 0, i, w, 1)
		glog.V(3).Infoln("Past: ", past)
		glog.V(3).Infoln("b: ", b)
		//The following performs this calculation
		//csSum += 1 - cos(2/pi*asin(norm(b - past*(past'*b))))
		blas.Gemv(
			past,
			b,
			v,
			matrix.FScalar(1.0),
			matrix.FScalar(0.0),
			linalg.OptTrans,
		)
		blas.Gemv(
			past,
			v,
			b,
			matrix.FScalar(-1.0),
			matrix.FScalar(1.0),
			linalg.OptNoTrans,
		)
		norm := blas.Nrm2(b).Float()
		csSum += (1 - math.Cos(math.Asin(norm)*2/math.Pi))
		glog.V(3).Infoln("cs: ", csSum)
	}
	return csSum / float64(lf)
}
