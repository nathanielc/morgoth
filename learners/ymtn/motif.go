package ymtn

import (
	"github.com/golang/glog"
	"github.com/nvcook42/matrix"
	"math"
	"math/rand"
	"sort"
)

var pInf = math.Inf(1)

type distanceFunc func(s1, s2 []float64) float64
type stopFunc func(x []float64, motifs, newMotifs []Motif) bool

type Motif struct {
	Beg int
	End int
}

func (self Motif) Diff() int {
	return self.End - self.Beg
}

var nWindows = 10

func DetectMotifs(x, changeScores []float64, min, max int) [][]Motif {
	locs := make([][]Motif, 0)

	var T = len(x)
	var extraLength = 0.2

	var candidateLocs = make([]int, nWindows)
	glog.V(1).Infoln("Raw CS:", changeScores)
	normalize(changeScores)
	glog.V(1).Infoln("NormCS:", changeScores)
	cumsum := calcCumSum(changeScores)
	glog.V(1).Infoln("CumSum:", cumsum)
	for i := 0; i < nWindows; i++ {
		p := rand.Float64()
		for j := 0; j < T; j++ {
			if p <= cumsum[j] {
				candidateLocs[i] = j
				break
			}
		}
	}

	w := int((1.0+extraLength)*float64(max) + 0.5)
	wbar := int(math.Max(3, math.Ceil(0.1*float64(w))))
	nSkip := int(math.Max(1, math.Floor(0.5*float64(wbar))))
	nSub := (w-wbar)/nSkip + 1
	nCompWindows := int(0.5 * float64(nWindows))
	dFun := dE
	stopFun := stopIfAboveMean

	iterations := int(math.Ceil(float64(nCompWindows) / float64(nWindows)))
	K := nWindows - 1

	glog.V(2).Infoln("w", w)
	glog.V(2).Infoln("wbar", wbar)
	glog.V(2).Infoln("nSkip", nSkip)
	glog.V(2).Infoln("nSub", nSub)
	glog.V(2).Infoln("nCompWindows", nCompWindows)
	glog.V(2).Infoln("iterations", iterations)
	glog.V(2).Infoln("K", K)

	//TODO make this dependent on the changeScores
	noise := generateNoiseWindow(w, x)

	allSW := createSWSet(candidateLocs, nWindows, nSub, w, T, nSkip)
	noiseSW := createSWSet([]int{int(math.Ceil(float64(w) / 2.0))}, 1, nSub, w, w, nSkip)

	glog.V(3).Infoln("noise:", noise)
	glog.V(3).Infoln("allSW:", allSW)
	glog.V(3).Infoln("noiseSW:", noiseSW)

	bestMatchesMn := compareFindBests(noise, x, allSW, noiseSW, wbar, K, dFun)
	glog.V(3).Infoln("bestMatchesMn:", bestMatchesMn)
	gamma := mean(bestMatchesMn)

	glog.V(3).Infoln("gamma:", gamma)

	inactiveSubs := make([]bool, nSub)
	for currentWin := 0; currentWin < nWindows; currentWin++ {
		glog.V(3).Infoln("currentWin", currentWin)
		for iteration := 0; iteration < iterations; iteration++ {
			candSW := make([]float64, 0, nSub)
			for i, inactiveSub := range inactiveSubs {
				if !inactiveSub {
					candSW = append(candSW, allSW.GetAt(currentWin, i))
				}
			}
			if len(candSW) == 0 {
				break
			}

			glog.V(3).Infoln("candSW", candSW)

			allSWSub := matrix.FloatZeros(nCompWindows, allSW.Cols())
			for i := 0; i < nCompWindows; i++ {
				row := rand.Intn(nWindows)
				for col := 0; col < allSWSub.Cols(); col++ {
					allSWSub.SetAt(i, col, allSW.GetAt(row, col))
				}
			}

			candSWMatrix := matrix.FloatVector(candSW).Transpose()
			bestMatchesM := compareFindBests(x, x, allSWSub, candSWMatrix, wbar, K, dFun)
			glog.V(3).Infoln("bestMatchesM:", bestMatchesM)

			updateCandSW(inactiveSubs, bestMatchesM, gamma)
		}
		anyActive := false
		for i := 0; i < len(inactiveSubs) && !anyActive; i++ {
			anyActive = !inactiveSubs[i]
		}
		if !anyActive {
			continue
		}

		mySubs := matrix.FloatZeros(1, allSW.Cols())
		allSW.SubMatrix(mySubs, currentWin, 0, 1, allSW.Cols())
		glog.V(3).Infoln("mySubs", mySubs)

		motifs := make([]Motif, 1)
		motifs[0] = Motif{-1, -1}
		nNotMatched := 0
		current := 0
		glog.V(3).Infoln("inactiveSubs", inactiveSubs)
		for i := 0; i < nSub; i++ {
			if motifs[current].Beg < 0 {
				if !inactiveSubs[i] {
					motifs[current].Beg = int(allSW.GetAt(currentWin, i))
					motifs[current].End = int(allSW.GetAt(currentWin, i))
				}
			} else {
				if !inactiveSubs[i] {
					motifs[current].End = int(allSW.GetAt(currentWin, i))
					nNotMatched = 0
				} else {
					nNotMatched += nSkip
					if nNotMatched >= wbar {
						motifs = append(motifs, Motif{-1, -1})
						current++
						nNotMatched = 0
					}
				}
			}
		}

		glog.V(3).Infoln("motifs", motifs)

		var maxMotif Motif
		max := -1
		for _, motif := range motifs {
			if motif.Beg >= 0 {
				diff := motif.Diff()
				if diff > max {
					max = diff
					maxMotif = motif
				}
			}
		}

		mean := x[maxMotif.Beg:maxMotif.End]

		defaultMotifs := detectMotifDefault(
			x,
			[]Motif{maxMotif},
			mean,
			candidateLocs,
			gamma,
			dFun,
			stopFun,
		)
		glog.V(2).Infoln("default motifs", defaultMotifs)
		if len(defaultMotifs) > 1 {
			locs = append(locs, defaultMotifs)
		}
	}

	return locs
}

func normalize(data []float64) {
	sum := 0.0
	for _, x := range data {
		sum += x
	}
	for i := 0; i < len(data); i++ {
		data[i] /= sum
	}
}

func calcCumSum(data []float64) []float64 {
	cumsum := make([]float64, len(data))
	sum := 0.0
	for i, x := range data {
		sum += x
		cumsum[i] = sum
	}
	for i := 0; i < len(data); i++ {
		data[i] /= sum
	}
	return cumsum
}

func generateNoiseWindow(w int, x []float64) []float64 {
	noise := make([]float64, w)
	size := len(x)
	for i := 0; i < w; i++ {
		noise[i] = x[rand.Intn(size)]
	}
	return noise
}

func createSWSet(candidateLocs []int, nWindows, nSub, w, T, nSkip int) *matrix.FloatMatrix {
	loc := matrix.FloatZeros(nWindows, nSub)
	glog.V(3).Infoln("createSWSet", nWindows, nSub)
	for i := 0; i < nWindows; i++ {
		be := int(
			math.Min(
				float64(T-w+1),
				math.Max(1, float64(candidateLocs[i])-math.Floor(float64(w)/2.0)),
			),
		)
		for j := 0; j < nSub; j++ {
			loc.SetAt(i, j, float64(be+j*nSkip))
		}
	}
	return loc
}

func dE(s1, s2 []float64) float64 {
	if len(s1) != len(s2) {
		panic("Cannot compute distance between different length series")
	}
	sum := 0.0
	for i := range s1 {
		diff := s1[i] - s2[i]
		sum += diff * diff
	}
	return sum / float64(len(s1))
}

func compareFindBests(
	cand []float64,
	comp []float64,
	allSW *matrix.FloatMatrix,
	candSW *matrix.FloatMatrix,
	wbar int,
	K int,
	dFun distanceFunc,
) []float64 {
	nWindows := allSW.Rows()
	nSub := allSW.Cols()
	nSubSrc := candSW.Cols()
	dists := matrix.FloatWithValue(nWindows, nSub, pInf)
	bestMatchesM := make([]float64, nSubSrc)
	glog.V(2).Infoln("compareFindBests", nWindows, nSub, nSubSrc)

	for k := 0; k < nSubSrc; k++ {
		matchDists := make([]float64, nWindows)
		for j := 0; j < nWindows; j++ {
			min := pInf
			for l := 0; l < nSub; l++ {
				compStart := int(allSW.GetAt(j, l))
				candStart := int(candSW.GetAt(0, k))
				s1 := comp[compStart : compStart+wbar-1]
				s2 := cand[candStart : candStart+wbar-1]

				distance := dFun(s1, s2)
				dists.SetAt(j, l, distance)

				if distance < min {
					min = distance
				}
			}
			matchDists[j] = min
		}
		sum := 0.0
		for i := range matchDists {
			sum += matchDists[i]
		}
		bestMatchesM[k] = sum / float64(nWindows)
	}

	sort.Float64s(bestMatchesM)
	return bestMatchesM
}

func mean(data []float64) float64 {
	sum := 0.0
	for _, x := range data {
		sum += x
	}
	return sum / float64(len(data))
}

func updateCandSW(inactiveSubs []bool, bestMatchesM []float64, gamma float64) {
	j := 0
	for k, inactiveSub := range inactiveSubs {
		if inactiveSub {
			continue
		}
		if bestMatchesM[j] > gamma {
			inactiveSubs[k] = true
		}
		j++
	}
}

func detectMotifDefault(
	x []float64,
	motifs []Motif,
	meanMotif []float64,
	candLocs []int,
	maxSmallDist float64,
	dFun distanceFunc,
	stopFun stopFunc,
) []Motif {

	glog.V(3).Infoln("detectMotifDefault", motifs)
	T := len(x)
	stepSize := 1
	length := len(meanMotif)

	nL := len(candLocs)
	locsSet := make(map[int]bool)
	for i := 0; i < nL; i++ {
		for j := -length; j < length+1; j += stepSize {
			loc := j + candLocs[i]
			if loc >= 0 && loc < T-length+1 {
				locsSet[loc] = true
			}
		}
	}

	nL = len(locsSet)
	locs := make([]int, nL)
	i := 0
	for loc := range locs {
		locs[i] = loc
		i++
	}

	glog.V(3).Infoln("locs:", locs)

	nB := len(motifs)
	toIgnore := make([]bool, nL, nL)

	var newMotifs []Motif

	for i = 0; i < nB; i++ {
		l := motifs[i].Diff() + 1
		for j := 0; j < nL; j++ {
			if locs[j] >= motifs[i].Beg-l+1 &&
				locs[j] <= motifs[i].End {
				toIgnore[j] = true
			}
		}
	}
	glog.V(3).Infoln("toIgnore:", toIgnore)
	for i = 0; i < nL; i++ {
		if toIgnore[i] {
			continue
		}
		beg, end := locs[i], locs[i]+length
		motif := x[beg:end]
		distance := dFun(meanMotif, motif)
		if distance > maxSmallDist {
			continue
		}
		newMotifs = append(motifs, Motif{beg, end})
		glog.V(3).Infoln("motifs:", motifs)
		glog.V(3).Infoln("newMotifs:", newMotifs)

		if len(motifs) > 1 {
			stop := stopFun(x, motifs, newMotifs)
			glog.V(3).Infoln("Stopping?", stop)
			if stop {
				continue
			}
		}
		motifs = newMotifs
		last := newMotifs[len(newMotifs)-1].End
		for i, loc := range locs {
			if loc < last {
				toIgnore[i] = true
			}
		}
	}
	return motifs
}

func stopIfAboveMean(x []float64, motifs, newMotifs []Motif) bool {

	glog.V(3).Infoln("stopIfAboveMean", motifs, newMotifs)

	epsilon := 1e-20
	acceptableMeanIncrease := 1.1

	n := len(motifs)
	m := maxDiff(motifs)

	X := matrix.FloatZeros(n, m)
	fillX(X, motifs, x)

	meanDistance := meanOfX(X)

	n = len(newMotifs)
	m = maxDiff(newMotifs)
	newX := matrix.FloatZeros(n, m)
	fillX(newX, newMotifs, x)

	newMeanDistance := meanOfX(newX)

	glog.V(3).Infoln("nmean, mean, max:", newMeanDistance, meanDistance, m)
	if newMeanDistance <= meanDistance*acceptableMeanIncrease {
		return false
	} else {
		return newMeanDistance-meanDistance > float64(m)*epsilon
	}
}

func maxDiff(motifs []Motif) int {
	max := -1
	for _, motif := range motifs {
		diff := motif.Diff()
		if diff > max {
			max = diff
		}
	}
	return max
}

func fillX(X *matrix.FloatMatrix, motifs []Motif, x []float64) {
	for j := 0; j < len(motifs); j++ {
		sub := x[motifs[j].Beg:motifs[j].End]
		for c := 0; c < len(sub); c++ {
			X.SetAt(j, c, sub[c])
		}
	}
}

func meanOfX(X *matrix.FloatMatrix) float64 {
	d := 0.0
	l := 0.0
	for i := 0; i < X.Rows()-1; i++ {
		for j := i + 1; j < X.Rows(); j++ {
			sum := 0.0
			for k := 0; k < X.Cols(); k++ {
				v := math.Abs(X.GetAt(i, k) - X.GetAt(j, k))
				sum += v
			}
			l += float64(X.Cols())
			d += sum
		}
	}
	return d / l
}
