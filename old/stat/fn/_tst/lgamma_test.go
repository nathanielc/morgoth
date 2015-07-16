package fn

import (
	"fmt"
	"testing"
)

// Test of ln(Gamma) against R:lgamma
func Test_LnΓ(t *testing.T) {
	const delta = 1e-4

	fmt.Println("Testing LnΓ #1")
	a := 3.15
	x := LnΓ(a)
	y := 0.8359236

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

	fmt.Println("Testing LnΓ #2")
	a = 432.123
	x = LnΓ(a)
	y = 2188.191

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

	fmt.Println("Testing LnΓ #3")
	a = 0.0000000675
	x = LnΓ(a)
	y = 16.51114

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

	fmt.Println("Testing LnΓ #4")
	a = 3.15e-19
	x = LnΓ(a)
	y = 42.60171

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

	fmt.Println("Testing LnΓ #5")
	a = 3.15e-99
	x = LnΓ(a)
	y = 226.8085

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

	fmt.Println("Testing LnΓ #6")
	a = -12.33
	x = LnΓ(a)
	y = -19.53042

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

}
