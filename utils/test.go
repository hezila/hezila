package utils

import (
	"fmt"
	"math"
	"testing"
)

func Expect(t *testing.T, expect string, actual interface{}) {
	actualString := fmt.Sprint(actual)
	if expect != actualString {
		t.Errorf("Expected value=\"%s\", Actual value=\"%s\"", expect, actualString)
	}
}

func ExpectNear(t *testing.T, expect float64, actual float64, acc float64) {
	if math.Abs(expect-actual) > acc {
		t.Errorf("Expected value=\"%s\", Actual value=\"%s\"", expect, actual)
	}
}
