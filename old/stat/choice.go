package stat

import (
	"math"
)

func Choice_PMF(θ []float64) func(i int64) float64 {
	return func(i int64) float64 {
		return θ[i]
	}
}
func Choice_LnPMF(θ []float64) func(i int64) float64 {
	return func(i int64) float64 {
		return log(θ[i])
	}
}
func NextChoice(θ []float64) int64 {
	u := NextUniform()
	i := 0
	sum := θ[0]
	for ; sum < u && i < len(θ)-1; i++ {
		sum += θ[i+1]
	}
	if u >= sum {
		return int64(len(θ))
	}
	return int64(i)
}
func Choice(θ []float64) func() int64 {
	return func() int64 {
		return NextChoice(θ)
	}
}
func NextLogChoice(lws []float64) int64 {
	return LogChoice(lws)()
}
func LogChoice(lws []float64) func() int64 {
	max := lws[0]
	for _, lw := range lws[1:len(lws)] {
		if lw > max {
			max = lw
		}
	}
	ws := make([]float64, len(lws))
	var sum float64
	for i, lw := range lws {
		ws[i] = math.Exp(lw - max)
		sum += ws[i]
	}
	norm := 1 / sum
	for i := range ws {
		ws[i] *= norm
	}
	return Choice(ws)
}
