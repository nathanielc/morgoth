package stat

func Bernoulli_PMF(ρ float64) func(k int64) float64 {
	return func(k int64) float64 {
		if k < 0 || k > 1 {
			panic("k is not 0 or 1")
		}
		if k == 1 {
			return ρ
		}
		return 1 - ρ
	}
}

func Bernoulli_PMF_At(ρ float64, k int64) float64 {
	pmf := Bernoulli_PMF(ρ)
	return pmf(k)
}

func Bernoulli_LnPMF(ρ float64) func(k int64) float64 {
	return func(k int64) float64 {
		if k == 1 {
			return log(ρ)
		}
		return log(1 - ρ)
	}
}

func NextBernoulli(ρ float64) int64 {
	if NextUniform() < ρ {
		return 1
	}
	return 0
}

func Bernoulli(ρ float64) func() int64 { return func() int64 { return NextBernoulli(ρ) } }

func Bernoulli_CDF(ρ float64) func(k int64) float64 {
	return func(k int64) float64 {
		if k < 0 || k > 1 {
			panic("k is not 0 or 1")
		}
		if k == 1 {
			return 1
		}
		return 1 - ρ
	}
}
