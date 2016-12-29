package math

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/*
Read-only matrix interface defines matrix operations that do not change the element value.
*/
type MatrixRO interface {
	// Return true if the matrix is nil.
	Nil() bool

	// Return the number of rows of this matrix
	Rows() uint

	// Return the number of columns of this matrix
	Cols() uint

	// Return the number of elements contained in this matrix
	NumElements() uint

	// Return the dimension of the matrix
	Dimension() (uint, uint)

	// Get the value in the ith row and jth column
	Get(i uint, j uint) float64

	// The determinant of this matrix
	Det() float64

	// The trace of this matrix
	Trace() float64

	// Returns an array of slices referencing the matrix data.
	// Changes to the slices effect changes to the matrix.
	Arrays() [][]float64

	// Returns the contents/slices of this matrix stored into a flat array (row-major).
	Array() []float64

	// The pretty-print string
	String() string

	SparseMatrix() *SparseMatrix
	DenseMatrix() *DenseMatrix
}

/*
A mutable matrix.
*/
type Matrix interface {
	MatrixRO

	// Set the value in the ith row and jth column
	Set(i uint, j uint, v float64)

	Add(MatrixRO) error
	Subtract(MatrixRO) error
	Scale(float64)
}

type matrix struct {
	rows uint
	cols uint
}

func (M *matrix) Nil() bool { return M == nil }

func (M *matrix) Rows() uint { return M.rows }

func (M *matrix) Cols() uint { return M.cols }

func (M *matrix) NumElements() uint { return M.rows * M.cols }

func (M *matrix) Dimension() (rows, cols uint) {
	rows = M.rows
	cols = M.cols
	return
}

/*
   Take a matlab-style matrix representation
   e.g., [a b c; d e f]
*/
func ParseMatlab(txt string) (A *DenseMatrix, err error) {
	var arrays [][]float64
	spaceSep := strings.Fields(txt)

	tok := func() (t string, eos bool) {
		defer func() {
			for len(spaceSep) != 0 && len(spaceSep[0]) == 0 {
				spaceSep = spaceSep[1:]
			}
		}()

		isNotNumber := func(c byte) bool {
			return c != '[' || c != ']' || c == ';'
		}

		if len(spaceSep) == 0 {
			eos = true
			return
		}

		top := spaceSep[0]

		var lof int
		for ; lof < len(top) && !isNotNumber(top[lof]); lof++ {
		}

		if lof != 0 {
			t = top[:lof]
			spaceSep[0] = top[lof:]
			return
		} else {
			t = top[:1]
			spaceSep[0] = top[1:]
			return
		}

		panic("unreadable")
	}

	stack := func(row []float64) (err error) {
		if len(arrays) == 0 {
			arrays = [][]float64{row}
			return
		}
		if len(arrays[0]) != len(row) {
			err = errors.New("misaligned row")
		}
		arrays = append(arrays, row)
		return
	}

	var row []float64

loop:
	for {
		t, eos := tok()
		if eos {
			break loop
		}
		switch t {
		case "[":
		case ";":
			err = stack(row)
			if err != nil {
				return
			}
			row = []float64{}
		case "]":
			err = stack(row)
			if err != nil {
				return
			}
			break loop
		default:
			var v float64
			v, err = strconv.ParseFloat(t, 64)
			if err != nil {
				return
			}
			row = append(row, v)
		}
	}

	A = MakeDenseMatrixStacked(arrays)
	return
}

func String(M MatrixRO) string {
	condense := func(vs string) string {
		if strings.Index(vs, ".") != -1 {
			for vs[len(vs)-1] == '0' {
				vs = vs[0 : len(vs)-1]
			}
		}
		if vs[len(vs)-1] == '.' {
			vs = vs[0 : len(vs)-1]
		}
		return vs
	}

	if M == nil {
		return "{nil}"
	}
	s := "{"

	var maxLen uint = 0
    var i, j uint
	for i = 0; i < M.Rows(); i++ {
		for j = 0; j < M.Cols(); j++ {
			v := M.Get(i, j)
			vs := condense(fmt.Sprintf("%f", v))

			maxLen = maxUInt(maxLen, uint(len(vs)))
		}
	}

	for i = 0; i < M.Rows(); i++ {
		for j = 0; j < M.Cols(); j++ {
			v := M.Get(i, j)

			vs := condense(fmt.Sprintf("%f", v))

			for uint(len(vs)) < maxLen {
				vs = " " + vs
			}
			s += vs
			if i != M.Rows()-1 || j != M.Cols()-1 {
				s += ","
			}
			if j != M.Cols()-1 {
				s += " "
			}
		}
		if i != M.Rows()-1 {
			s += "\n "
		}
	}
	s += "}"
	return s
}
