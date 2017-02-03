package math

import (
	"math/rand"

	log "github.com/golang/glog"
)

// A matrix backed by a flat array of all elements
type DenseMatrix struct {
	matrix

	// flatted elements
	elements []float64

	// offset between rows; step = cols
	step uint
}

func NewDenseMatrix(rows, cols uint) *DenseMatrix {
	M := new(DenseMatrix)
	M.rows = rows
	M.cols = cols
	M.step = cols

	M.elements = make([]float64, rows*cols)

	for i := uint(0); i < rows*cols; i++ {
		M.elements[i] = 0.0
	}
	return M
}

func MakeDenseMatrix(elements []float64, rows, cols uint) *DenseMatrix {
	A := new(DenseMatrix)
	A.rows = rows
	A.cols = cols
	A.step = cols

	A.elements = elements
	return A
}

func MakeDenseMatrixStacked(data [][]float64) *DenseMatrix {
	rows := uint(len(data))
	cols := uint(len(data[0]))
	elements := make([]float64, rows*cols)

	for i := uint(0); i < rows; i++ {
		for j := uint(0); j < cols; j++ {
			elements[i*cols+j] = data[i][j]
		}
	}
	return MakeDenseMatrix(elements, rows, cols)
}

func (M *DenseMatrix) Arrays() [][]float64 {
	a := make([][]float64, M.rows)
	for i := uint(0); i < M.rows; i++ {
		a[i] = M.elements[i*M.step : i*M.step+M.cols]
	}
	return a
}

func (M *DenseMatrix) Array() []float64 {
	if M.step == M.cols {
		return M.elements[0 : M.rows*M.cols]
	}

	a := make([]float64, M.rows*M.cols)
	for i := uint(0); i < M.rows; i++ {
		for j := uint(0); j < M.cols; j++ {
			a[i*M.cols+j] = M.elements[i*M.step+j]
		}
	}
	return a
}

func (M *DenseMatrix) RowSlice(row uint) []float64 {
	if row < 0 || row > M.rows-1 {
		log.Fatal("index out of bound!")
	}
	return M.elements[row*M.step : row*M.step+M.cols]
}

func (M *DenseMatrix) ColSlice(col uint) []float64 {
	if col < 0 || col > M.cols-1 {
		log.Fatal("index out of bound!")
	}
	a := make([]float64, M.rows)
	for i := uint(0); i < M.rows; i++ {
		a[i] = M.elements[i*M.step+col]
	}
	return a
}

func (M *DenseMatrix) Get(i, j uint) (v float64) {
	if i < 0 {
		i = M.rows + i
		if i < 0 {
			log.Fatal("index out of bound!")
		}
	}

	if j < 0 {
		j = M.cols + j
		if j < 0 {
			log.Fatal("index out of bound!")
		}
	}

	if i >= M.rows || j >= M.cols {
		log.Fatal("index out of bound!")
	}
	v = M.elements[i*M.step+j]
	return
}

func (M *DenseMatrix) Set(i, j uint, v float64) {
	if i < 0 {
		i = M.rows + i
		if i < 0 {
			log.Fatal("index out of bound!")
		}
	}
	if j < 0 {
		j = M.cols + j
		if j < 0 {
			log.Fatal("index out of bound!")
		}
	}
	if i >= M.rows || j >= M.cols {
		log.Fatal("index out of bound!")
	}
	M.elements[i*M.step+j] = v
	return
}

func (M *DenseMatrix) SetValue(index uint, v float64) {
	M.elements[index] = v
}

// Get a submatrix starting at i, j with rows rows and cols columns
func (M *DenseMatrix) SubMatrix(i, j, rows, cols uint) (A *DenseMatrix, err error) {
	if i < 0 || j < 0 || rows <= 0 || cols <= 0 ||
		(i+rows) > M.rows || (j+cols) > M.cols {
		err = ErrorIllegalIndex
	}

	A = new(DenseMatrix)
	A.elements = make([]float64, rows*cols)
	A.step = cols
	A.rows = rows
	A.cols = cols
	for r := uint(0); r < rows; r++ {
		for c := uint(0); c < cols; c++ {
			A.elements[r*A.step+c] = M.elements[(i+r)*M.step+j+c]
		}
	}
	return
}

// Get a submatrix starting at i,j with rows rows and cols columns. Changes to
// the returned matrix show up in the original.
func (A *DenseMatrix) GetMatrix(i, j, rows, cols uint) *DenseMatrix {
	B := new(DenseMatrix)
	B.elements = A.elements[i*A.step+j : i*A.step+j+(rows-1)*A.step+cols]
	B.rows = rows
	B.cols = cols
	B.step = A.step
	return B
}

// Copy A into M, with A's 0, 0 aligning with M's i, j
func (M *DenseMatrix) SetMatrix(i, j uint, A *DenseMatrix) {
	var r, c uint
	for r := uint(0); r < A.rows; r++ {
		for c := uint(0); c < A.cols; c++ {
			M.Set(i+r, j+c, A.Get(r, c))
		}
	}
}

func (M *DenseMatrix) ColVector(j uint) *DenseMatrix {
	return M.SubMatrix(0, j, M.rows, 1)
}

func (M *DenseMatrix) RowVector(i uint) *DenseMatrix {
	return M.SubMatrix(i, 0, 1, M.cols)
}

func (M *DenseMatrix) Copy() *DenseMatrix {
	A := new(DenseMatrix)
	A.rows = M.rows
	A.cols = M.cols
	A.step = M.step
	A.elements = make([]float64, M.rows*M.cols)

	for r := uint(0); r < A.rows; r++ {
		copy(A.RowSlice(r), M.RowSlice(r))
	}
	return A
}

// Get a new matrix [M, A], with A after M
func (M *DenseMatrix) Augment(A *DenseMatrix) (B *DenseMatrix, err error) {
	if M.rows != A.rows {
		err = ErrorDimensionMismatch
		return
	}
	B = Zeros(M.rows, M.cols+A.cols)
	err = M.AugmentFill(A, B)
	return
}

// Get a new matrix [M; A], with M above A
func (M *DenseMatrix) Stack(A *DenseMatrix) (B *DenseMatrix, err error) {
	if M.cols != A.cols {
		err = ErrorDimensionMismatch
		return
	}
	B = Zeros(M.rows+A.rows, M.cols)
	err = M.StackFill(A, B)
	return
}

func (M *DenseMatrix) StackFill(A, B *DenseMatrix) (err error) {
	if M.cols != A.cols || M.cols != B.cols || B.rows != M.rows+A.rows {
		err = ErrorDimensionMismatch
		return
	}
	B.SetMatrix(0, 0, M)
	B.SetMatrix(M.rows, 0, A)
	return
}

func (M *DenseMatrix) AugmentFill(A, B *DenseMatrix) (err error) {
	if M.rows != A.rows || M.rows != B.rows || B.cols != M.rows+A.rows {
		err = ErrorDimensionMismatch
		return
	}

	B.SetMatrix(0, 0, M)
	B.SetMatrix(0, A.cols, A)
	return
}

// Create a sparse matrix copy.
func (A *DenseMatrix) SparseMatrix() *SparseMatrix {
	B := ZerosSparse(A.rows, A.cols)
	for i := uint(0); i < A.rows; i++ {
		for j = uint(0); j < A.cols; j++ {
			v := A.Get(i, j)
			if v != 0 {
				B.Set(i, j, v)
			}
		}
	}
	return B
}

func (A *DenseMatrix) DenseMatrix() *DenseMatrix {
	return A.Copy()
}

func Zeros(rows, cols uint) *DenseMatrix {
	Z := new(DenseMatrix)
	Z.elements = make([]float64, rows*cols)
	Z.rows = rows
	Z.cols = cols
	Z.step = cols
	return Z
}

func Ones(rows, cols uint) *DenseMatrix {
	O := new(DenseMatrix)
	O.elements = make([]float64, rows*cols)
	O.rows = rows
	O.cols = cols
	O.step = cols
	for i := uint(0); i < rows*cols; i++ {
		O.elements[i] = 1
	}
	return O
}

func Eye(size uint) *DenseMatrix {
	E := Zeros(size, size)

	for i := uint(0); i < size; i++ {
		E.Set(i, i, 1)
	}
	return E
}

func Normals(rows, cols uint) *DenseMatrix {
	N := Zeros(rows, cols)

	for i := uint(0); i < N.Rows(); i++ {
		for j := uint(0); j < N.Cols(); j++ {
			N.Set(i, j, rand.NormFloat64())
		}
	}
	return N
}

func Diagonal(d []float64) *DenseMatrix {
	n := uint(len(d))
	D := Zeros(n, n)

	for i := uint(0); i < n; i++ {
		D.Set(i, i, d[i])
	}
	return D
}

func MakeDenseCopy(A MatrixRO) *DenseMatrix {
	B := Zeros(A.Rows(), A.Cols())

	for i := uint(0); i < B.rows; i++ {
		for j := uint(0); j < B.cols; j++ {
			B.Set(i, j, A.Get(i, j))
		}
	}
	return B
}

func (A *DenseMatrix) String() string { return String(A) }
