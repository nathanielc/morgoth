package stat

import (
	. "github.com/nathanielc/morgoth/stat/fn"
)

func StudentsT_PDF(ν float64) func(x float64) float64 {
	normalization := Γ((ν+1)/2) / (sqrt(ν*π) * Γ(ν/2))
	return func(x float64) float64 {
		return normalization * pow(1+x*x/ν, -(ν+1)/2)
	}
}
func StudentsT_LnPDF(ν float64) func(x float64) float64 {
	normalization := LnΓ((ν+1)/2) - log(sqrt(ν*π)) - LnΓ(ν/2)
	return func(x float64) float64 {
		return normalization + log(1+x*x/ν)*-(ν+1)/2
	}
}

//StudentsT(ν) => N(0, 1)*sqrt(ν/NextGamma(ν/2, 2))
func NextStudentsT(ν float64) float64 {
	return NextNormal(0, 1) * sqrt(ν/NextGamma(ν/2, 2))
}
func StudentsT(ν float64) func() float64 {
	return func() float64 {
		return NextStudentsT(ν)
	}
}
