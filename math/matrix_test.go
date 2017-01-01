package math

import (
	"testing"

	"github.com/hezila/hezila/utils"
)

func TestParse(t *testing.T) {
	s := `[1 2 3; 4 5 6]`
	A, err := ParseMatlab(s)

	if err != nil {
		t.Fatal(err)
	}

	Ar := MakeDenseMatrix([]float64{1, 2, 3, 4, 5, 6}, 2, 3)
	if !Equals(A, Ar) {
		t.Error()
	}
}

func TestString(t *testing.T) {
	A := MakeDenseMatrix([]float64{1, 2, 3, 4, 5, 6}, 2, 3)
	s := `{1, 2, 3 \n 4, 5, 6}`
	utils.Expect(t, s, String(A))
}
