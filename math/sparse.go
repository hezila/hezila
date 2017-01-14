package math

import (
	"log"
	"math/rand"

	"hezila/utils"
)

// A sparse matrix with indexing all of its elements by a map
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

func NewSparseMatrix(rows, cols uint) *SparseMatrix {
	M := new(SparseMatrix)
	M.rows = rows
	M.cols = cols
	M.offset = 0
	M.step = cols
	
	M.elements = make(map[uint]float64)
	return M
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
	log.Fatal("The Arrays() function for sparse matrix has not been implemented!")
	return nil
}

// TODO: to implements
func (M *SparseMatrix) Array() []float64 {
	log.Fatal("The Array() function for sparse matrix has not been implemented!")
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

func (M *SparseMatrix) Get(i, j uint) (float64) {
	if i < 0 {
		i = M.rows + i
		if i < 0 {
			log.Fatal("index out of bound!")
			//err = ErrorIllegalIndex
			return nil
		}
	}

	if j < 0 {
		j = M.cols + j
		if j < 0 {
			log.Fatal("index out of bound!")
			//err = ErrorIllegalIndex
			return nil
		}
	}

	if i >= M.rows || j >= M.cols {
		log.Fatal("index out of bound!")
		// err = ErrorIllegalIndex
		return nil
	}
	
	v, err = M.elements[i*M.step+j+M.offset]
	if err != nil {
		log.Fatal("the element indexed does not exists!")
	}
	return v
}

func (M *SparseMatrix) Exist(i, j uint) (v float64, err error) {
	if i < 0 {
		i = M.rows + i
		if i < 0 {
			err = ErrorIllegalIndex
		}
	}

	if j < 0 {
		j = M.cols + j
		if j < 0 {
			err = ErrorIllegalIndex
		}
	}

	if i >= M.rows || j >= M.cols {
		err = ErrorIllegalIndex
	}
	
	v, err = M.elements[i*M.step+j+M.offset]
	return
}

// Looks up an element given its element index
func (M *SparseMatrix) GetValue(index uint) (v float64, err error) {
	v, err = M.elements[index]
	return
}

func (M *SparseMatrix) Set(i, j uint, v float64) {
	//i = i % M.rows
	if i < 0 {
		i = M.rows + i
		if i < 0 {
			log.Fatal("index out of bound!")
			//err = ErrorIllegalIndex
		}
	}

	//j = j % M.cols
	if j < 0 {
		j = M.cols + j
		if j < 0 {
			log.Fatal("index out of bound!")
			//err = ErrorIllegalIndex
		}
	}

	if v == nil {
		delete(M.elements, i*M.step+j+M.offset)
	} else {
		M.elements[i*M.step+j+M.offset] = v
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
			if 0 <= i && i < M.rows && 0 <= j && j < M.cols {
				o <- index
			}
		}
		close(o)
	}(out)
	return
}

func (M *SparseMatrix) SubMatrix(i, j, rows, cols uint) (*SparseMatrix, err error) {
	if i < 0 || j < 0 || rows <= 0 || cols <= 0 ||
		(i+rows) > M.rows || (j+cols) > M.cols {
		err = ErrorIllegalIndex
	}

	S := ZerosSparse(rows, cols)
	for r := uint(0); r < rows; r++ {
		for c := uint(0); c < cols; c++ {
			index = (i+r)*M.step + (j+c) + M.offset
			if val, ok := M.elements[index]; ok {
				S.Set(r, c, val)
			}
		}
	}

	return (S, err)
}

func (M *SparseMatrix) ColVector(j uint) *SparseMatrix {
	return M.SubMatrix(0, j, M.rows, 1)
}

func (M *SparseMatrix) RowVector(i uint) *SparseMatrix {
	return M.SubMatrix(i, 0, 1, M.cols)
}

// Create a new matrix [A B]
func (A *SparseMatrix) Augment(B *SparseMatrix) (S *SparseMatrix, err error) {
	if A.rows != B.rows {
		err = ErrorDimensionMismatch
		return
	}

	S = ZerosSparse(A.rows, A.cols+B.cols)

	for index, value := range A.elements {
		i, j := A.GetRowColIndex(index)
		S.Set(i, j, value)
	}

	for index, value := range B.elements {
		i, j := B.GetRowColIndex(index)
		S.Set(i, j+A.cols, value)
	}

	return
}

func (A *SparseMatrix) Stack(B *SparseMatrix) (S *SparseMatrix, err error) {
	if A.cols != B.cols {
		err = ErrorDimensionMismatch
		return
	}

	S = ZerosSparse(A.rows+B.rows, A.cols)

	for index, value := range A.elements {
		i, j := A.GetRowColIndex(index)
		S.Set(i, j, value)
	}

	for index, value := range B.elements {
		i, j := B.GetRowColIndex(index)
		S.Set(i+A.rows, j, value)
	}

	return
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
	return M
}

func OnesSparse(rows, cols uint) *SparseMatrix {
	O := new(SparseMatrix)
	O.rows = rows
	O.cols = cols
	O.step = cols
	O.elements = map[uint]float64{}

	for i := uint(0); i < cols*cols; i++ {
		O.elements[i] = 1
	}
	return O
}

func EyeSparse(size uint) *SparseMatrix {
	E := ZerosSparse(size, size)

	for i := uint(0); i < size; i++ {
		E.Set(i, i, 1)
	}
	return E
}

func NormalsSparse(rows, cols uint) *SparseMatrix {
	N := ZerosSparse(rows, cols)

	for i := uint(0); i < rows; i++ {
		for j := uint(0); j < cols; j++ {
			N.Set(i, j, rand.NormFloat64())
		}
	}
	return N
}

func Diagonal(d []float64) *SparseMatrix {
	n := uint(len(d))
	D := ZerosSparse(n, n)
	for i := uint(0); i < n; i++ {
		D.Set(i, i, d[i])
	}
	return D
}

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

func MakeSparseCopy(M MatrixRO) *SparseMatrix {
	A := ZerosSparse(M.Rows(), M.Cols())

	for i := uint(0); i < M.Rows(); i++ {
		for j := uint(0); j < M.Cols(); j++ {
			A.Set(i, j, M.Get(i, j))
		}
	}
	return A
}

func (A *SparseMatrix) String() string { return String(A) }
