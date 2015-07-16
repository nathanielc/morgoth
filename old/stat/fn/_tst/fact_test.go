package fn

import (
	"fmt"
	"testing"
)

// Test against R:factorial
func Test_Fact(t *testing.T) {

	fmt.Println("Testing Fact #1")
	delta := 1e152
	var n int64 = 100
	x := Fact(n)
	y := 9.332622e+157

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

	fmt.Println("Testing Fact #2")
	delta = 1e175
	n = 111
	x = Fact(n)
	y = 1.762953e+180

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

	fmt.Println("Testing Fact #3")
	delta = 1e301
	n = 170
	x = Fact(n)
	y = 7.257416e+306

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

}

// Test against R:lfactorial
func Test_LnFact(t *testing.T) {

	fmt.Println("Testing LnFact #1")
	delta := 1e-3
	var n float64 = 170
	x := LnFact(n)
	y := 706.5731

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

	fmt.Println("Testing LnFact #2")
	delta = 1
	n = 300000
	x = LnFact(n)
	y = 3483469

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

	fmt.Println("Testing LnFact #3")
	delta = 1
	n = 3584549
	x = LnFact(n)
	y = 50513986

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

	fmt.Println("Testing LnFact #4")
	delta = 1e9
	n = 3584549335689
	x = LnFact(n)
	y = 1.000364e+14

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

	fmt.Println("Testing LnFact #5")
	delta = 1e36
	n = 3584549335689779664547856867589564784475
	x = LnFact(n)
	y = 3.228871e+41

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

	fmt.Println("Testing LnFact #6")
	delta = 1e55
	n = 35845493356897796645478568675895647844758548848745565856895
	x = LnFact(n)
	y = 4.797079e+60

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

	fmt.Println("Testing LnFact #7")
	delta = 1e298
	n = 3.695685e300
	x = LnFact(n)
	y = 2.554024e+303

	if abs(x-y) > delta {
		fmt.Println("failed: ", x, y)
		t.Error()
	}

}
