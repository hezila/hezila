package math

import (
	"math/rand"
)

type SparseMatrix struct {
	matrix

	// flatted elements
	elements map[uint]float64

	// offset to start of matrix s.t. idx = i * cols + j + offset
	// offset  = staring row * step + staring col
	offset uint
	// analogous to dense step
	step uint
}

func MakeSparseMatrix(elements map[uint]float64, rows, cols uint) *SparseMatrix {
	M := new(SparseMatrix)
	M.rows = rows
	M.cols = cols
	M.offset = 0
	M.step = cols

	M.elements = elements
	return M
}

// TODO: to implements
func (M *SparseMatrix) Arrays() [][]float64 {
	return nil
}

// TODO: to implements
func (M *SparseMatrix) Array() []float64 {
	return nil
}

func (M *SparseMatrix) GetRowColIndex(index uint) (i, j uint) {
	i = (index - M.offset) / M.step
	j = (index - M.offset) % M.step
	return
}

func (M *SparseMatrix) GetRowIndex(index uint) (i uint) {
	i = (index - M.offset) / M.step
	return
}

func (M *SparseMatrix) GetColIndex(index uint) (j uint) {
	j = (index - M.offset) % M.step
	return
}

func (M *SparseMatrix) Get(i, j uint) float64 {
	i = i % M.rows
	if i < 0 {
		i = M.rows + i
	}

	j = j % M.cols
	if j < 0 {
		j = M.cols + j
	}

	v, ok := M.elements[i*M.step+j+M.offset]
	if !ok {
		return nil
	}
	return v
}

// Looks up an element given its element index
func (M *SparseMatrix) GetValue(index uint) float64 {
	v, ok := M.elements[index]
	if !ok {
		return nil
	}
	return v
}

func (M *SparseMatrix) Set(i, j uint, v float64) {
	i = i % M.rows
	if i < 0 {
		i = M.rows + i
	}

	j = j % M.cols
	if j < 0 {
		j = M.cols + j
	}
	index = i*M.step + j + M.offset
	if v == nil {
		delete(M.elements, index)
	} else {
		M.elements[index] = v
	}

}

func (M *SparseMatrix) SetValue(index uint, v float64) {
	if v == nil {
		delete(M.elements, index)
	} else {
		M.elements[index] = v
	}
}

func (M *SparseMatrix) Indices() (out chan uint) {
	// maybe thread the populating?
	out = make(chan uint)
	go func(o chan uint) {
		for index := range M.elements {
			i, j := M.GetRowColIndex(index)
			if i >= 0 && i < M.rows && j >= 0 && j < M.cols {
				o <- index
			}
		}
		close(o)
	}(out)
	return
}

func (M *SparseMatrix) SubMatrix(i, j, rows, cols uint) *SparseMatrix {
	if i < 0 || j < 0 || i+rows > M.rows || j+cols > M.cols {
		i = maxUInt(0, i)
		j = maxUInt(0, j)
		rows = minUInt(M.rows-i, rows)
		cols = minUInt(M.cols-j, cols)
	}
	S := ZerosSparse(rows, cols)

	for index, value := range M.elements {
		r, c := M.GetRowColIndex(index)
		if r < i+rows && c < j+cols {
			S.Set(r-i, c-j, value)
		}
	}
	return S
}

func (M *SparseMatrix) ColVector(j uint) *SparseMatrix {
	return M.SubMatrix(0, j, M.rows, 1)
}

func (M *SparseMatrix) RowVector(i uint) *SparseMatrix {
	return M.SubMatrix(i, 0, 1, M.cols)
}

// Create a new matrix [A B]
func (A *SparseMatrix) Augment(B *SparseMatrix) (*SparseMatrix, error) {
	if A.rows != B.rows {
		return nil, ErrorDimensionMismatch
	}

	C := ZerosSparse(A.rows, A.cols+B.cols)

	for index, value := range A.elements {
		i, j := A.GetRowColIndex(index)
		C.Set(i, j, value)
	}

	for index, value := range B.elements {
		i, j := B.GetRowColIndex(index)
		C.Set(i, j+A.cols, value)
	}

	return C, nil
}

func (A *SparseMatrix) Stack(B *SparseMatrix) (*SparseMatrix, error) {
	if A.cols != B.cols {
		return nil, ErrorDimensionMismatch
	}

	C := ZerosSparse(A.rows+B.rows, A.cols)

	for index, value := range A.elements {
		i, j := A.GetRowColIndex(index)
		C.Set(i, j, value)
	}

	for index, value := range B.elements {
		i, j := B.GetRowColIndex(index)
		C.Set(i+A.rows, j, value)
	}

	return C, nil
}

func (M *SparseMatrix) L() *SparseMatrix {
	B := ZerosSparse(M.rows, M.cols)
	for index, value := range M.elements {
		i, j := M.GetRowColIndex(index)
		if i >= j {
			B.Set(i, j, value)
		}
	}
	return B
}

func (M *SparseMatrix) U() *SparseMatrix {
	U := ZerosSparse(M.rows, M.cols)
	for index, value := range M.elements {
		i, j := M.GetRowColIndex(index)
		if i <= j {
			U.Set(i, j, value)
		}
	}
	return U
}

func (M *SparseMatrix) Copy() *SparseMatrix {
	C := ZerosSparse(M.rows, M.cols)
	for index, value := range M.elements {
		C.elements[index] = value
	}
	return C
}

func ZerosSparse(rows, cols uint) *SparseMatrix {
	M := new(SparseMatrix)
	M.rows = rows
	M.cols = cols
	M.offset = 0
	M.step = cols
	M.elements = map[uint]float64{}
	var i uint
	for i = 0; i < rows*cols; i++ {
		M.elements[i] = 0
	}
	return M
}

func OnesSparse(rows, cols uint) *SparseMatrix {
	O := new(SparseMatrix)
	O.rows = rows
	O.cols = cols
	O.step = cols
	O.elements = map[uint]float64{}
	var i uint
	for i = 0; i < cols*cols; i++ {
		O.elements[i] = 1.0
	}
	return O
}

func EyeSparse(size uint) *SparseMatrix {
	E := ZerosSparse(size, size)
	var i uint
	for i = 0; i < size; i++ {
		E.Set(i, i, 1.0)
	}
	return E
}

func NormalsSparse(rows, cols uint) *SparseMatrix {
	N := ZerosSparse(rows, cols)
	var i, j uint
	for i = 0; i < rows; i++ {
		for j = 0; j < cols; j++ {
			N.Set(i, j, rand.NormFloat64())
		}
	}
	return N
}

//func Diagonal(d []float64) *SparseMatrix {
//	n := len(d)
//	D := ZerosSparse(n, n)
//	for i := 0; i < n; i++ {
//		D.Set(i, i, d[i])
//	}
//	return D
//}

/*
Convert this sparse matrix into a dense matrix.
*/
func (A *SparseMatrix) DenseMatrix() *DenseMatrix {
	B := Zeros(A.rows, A.cols)
	for index, value := range A.elements {
		i, j := A.GetRowColIndex(index)
		B.Set(i, j, value)
	}
	return B
}

func (A *SparseMatrix) SparseMatrix() *SparseMatrix {
	return A.Copy()
}

func (A *SparseMatrix) String() string { return String(A) }

func MakeSparseCopy(M Matrix) *SparseMatrix {
	A := ZerosSparse(M.Rows(), M.Cols())
	var i, j uint
	for i = 0; i < M.Rows(); i++ {
		for j = 0; j < M.Cols(); j++ {
			A.Set(i, j, M.Get(i, j))
		}
	}
	return A
}
