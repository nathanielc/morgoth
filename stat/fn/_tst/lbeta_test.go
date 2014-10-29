package fn

import (
	"fmt"
	"testing"
)

// Test against R:lbeta
func Test_LnBeta(t *testing.T) {
	const delta = 1e-4

	fmt.Println("Testing LnBeta #1")
	a, b := 3.15, 4.22
	x := LnBeta(a, b)
	y := -4.371739

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

	fmt.Println("Testing LnBeta #2")
	a, b = 215.666, 0.333
	x = LnBeta(a, b)
	y = -0.8024721

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

	fmt.Println("Testing LnBeta #3")
	a, b = 0.000123, 22.334
	x = LnBeta(a, b)
	y = 9.002876

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

}
