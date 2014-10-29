package stat

func Geometric_PMF(ρ float64) func(i int64) float64 {
	return func(n int64) float64 { return ρ * pow(ρ, float64(n)) }
}
func Geometric_LnPMF(ρ float64) func(i int64) float64 {
	return func(n int64) float64 { return log(ρ) + float64(n)*log(ρ) }
}

//NextGeometric(ρ) => # of NextBernoulli(ρ) failures before one success
func NextGeometric(ρ float64) int64 {
	if NextBernoulli(ρ) == 1 {
		return 1 + NextGeometric(ρ)
	}
	return 0
}
func Geometric(ρ float64) func() int64 { return func() int64 { return NextGeometric(ρ) } }
