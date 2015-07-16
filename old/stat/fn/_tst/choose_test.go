package fn

import (
	"fmt"
	"testing"
)

// Test against R:choose

func TestFChoose(t *testing.T) {
	const delta = 1e-4

	fmt.Println("Testing FChoose #1")
	n, k := 15.3, 4.0
	a := FChoose(n, k)
	x := 1491.327

	if abs(x-a) > delta {
		fmt.Println("failed: ", a, x)
		t.Error()
	}

	fmt.Println("Testing FChoose #2")
	n, k = 21.3234, 6.0
	a = FChoose(n, k)
	x = 60263.41

	if abs(x-a) > delta {
		fmt.Println("failed: ", a, x)
		t.Error()
	}

	fmt.Println("Testing FChoose #3")
	n, k = 21.3234, 34.0
	a = FChoose(n, k)
	x = 2.690386e-11

	if abs(x-a) > delta {
		fmt.Println("failed: ", a, x)
		t.Error()
	}

	fmt.Println("Testing FChoose #4")
	n, k = 0.77, 34.0
	a = FChoose(n, k)
	x = -0.0003862904

	if abs(x-a) > delta {
		fmt.Println("failed: ", a, x)
		t.Error()
	}
}

func TestLnFChoose(t *testing.T) {
	const delta = 1e-4

	fmt.Println("Testing LnFChoose #1")
	n, k := 15.3, 4.0
	a := LnFChoose(n, k)
	x := 7.307422

	if abs(x-a) > delta {
		fmt.Println("failed: ", a, x)
		t.Error()
	}

	fmt.Println("Testing LnFChoose #2")
	n, k = 21.3234, 6.0
	a = LnFChoose(n, k)
	x = 11.00648

	if abs(x-a) > delta {
		fmt.Println("failed: ", a, x)
		t.Error()
	}

	fmt.Println("Testing LnFChoose #3")
	n, k = 21.3234, 34.0
	a = LnFChoose(n, k)
	x = -24.33875

	if abs(x-a) > delta {
		fmt.Println("failed: ", a, x)
		t.Error()
	}

	fmt.Println("Testing LnFChoose #4")
	n, k = 0.77, 34.0
	a = LnFChoose(n, k)
	x = -7.858921

	if abs(x-a) > delta {
		fmt.Println("failed: ", a, x)
		t.Error()
	}
}
