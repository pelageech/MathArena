package math

import (
	"testing"
)

func TestCalculate(t *testing.T) {
	sum := Sum{
		Num(23), Num(10), Num(1),
	}

	if sum.Calculate() != 34 {
		t.Fatalf("Calculate fail: got %d, want %d", sum.Calculate(), 34)
	}
}

func TestExpressionJSON(t *testing.T) {
	sum := Sum{
		Num(23), Sum{
			Num(10), Num(90), Num(-9),
		}, Num(1),
	}

	if string(sum.Marshal()) != "23+10+90-9+1" {
		t.Fatalf("expected: `%v`, got: `%v`", "23+10+90-9+1", string(sum.Marshal()))
	}
}
