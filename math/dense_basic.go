package math

import (
	"errors"
	"math"
)

func (A *DenseMatrix) Symmetric() bool {
	if A.rows != A.cols {
		return false
	}
	for i := uint(0); i < A.rows; i++ {
		for j := uint(0); j < i; j++ {
			if A.Get(i, j) != A.Get(j, i) {
				return false
			}
		}
	}
	return true
}

func (M *DenseMatrix) SwapRows(r1, r2 uint) {
	index1 := r1 * M.step
	index2 := r2 * M.step

	for j := uint(0); j < M.cols; j++ {
		temp := M.elements[index1]
		M.elements[index1] = M.elements[index2]
		M.elements[index2] = temp
		index1++
		index2++
	}
}

func (M *DenseMatrix) ScaleRow(r uint, f float64) {
	index := r * M.step

	for j := uint(0); j < M.cols; j++ {
		M.elements[index] *= f
		index++
	}
}

func (M *DenseMatrix) ScaleAddRow(rd, rs uint, f float64) {
	indexd := rd * M.step
	indexs := rs * M.step

	for j := uint(0); j < M.cols; j++ {
		M.elements[indexd] += f * M.elements[indexs]
		indexd++
		indexs++
	}
}

func (A *DenseMatrix) Inverse() (*DenseMatrix, error) {
	if A.rows != A.cols {
		return nil, ErrorDimensionMismatch
	}
	aug, _ := A.Augment(Eye(A.rows))

	for i := uint(0); i < aug.rows; i++ {
		j := i
		for k := i; k < aug.rows; k++ {
			if math.Abs(aug.Get(k, i)) > math.Abs(aug.Get(j, i)) {
				j = k
			}
		}
		if j != i {
			aug.SwapRows(i, j)
		}
		if aug.Get(i, i) == 0 {
			return nil, ExceptionSingular
		}
		aug.ScaleRow(i, 1.0/aug.Get(i, i))

		for k := uint(0); k < aug.rows; k++ {
			if k == i {
				continue
			}
			aug.ScaleAddRow(k, i, -aug.Get(k, i))
		}
	}
	inv := aug.GetMatrix(0, A.cols, A.rows, A.cols)
	return inv, nil
}

func (A *DenseMatrix) Det() float64 {
	B := A.Copy()
	P := B.LUInPlace()
	return product(B.DiagonalCopy()) * P.Det()
}

func (A *DenseMatrix) Trace() float64 { return sum(A.DiagonalCopy()) }

func (A *DenseMatrix) OneNorm() (ε float64) {
	var i, j uint
	for i = 0; i < A.rows; i++ {
		for j = 0; j < A.cols; j++ {
			ε = max(ε, A.Get(i, j))
		}
	}
	return
}

func (A *DenseMatrix) TwoNorm() float64 {
	var sum float64 = 0
	var i, j uint
	for i = 0; i < A.rows; i++ {
		for j = 0; j < A.cols; j++ {
			v := A.elements[i*A.step+j]
			sum += v * v
		}
	}
	return math.Sqrt(sum)
}

func (A *DenseMatrix) InfinityNorm() (e float64) {
	var i, j uint
	for i = 0; i < A.rows; i++ {
		for j = 0; j < A.cols; j++ {
			e += A.Get(i, j)
		}
	}
	return
}

func (A *DenseMatrix) Transpose() *DenseMatrix {
	B := Zeros(A.Cols(), A.Rows())
	var i, j uint
	for i = 0; i < A.Rows(); i++ {
		for j = 0; j < A.Cols(); j++ {
			B.Set(j, i, A.Get(i, j))
		}
	}
	return B
}

func (A *DenseMatrix) TransposeInPlace() (err error) {
	if A.rows != A.cols {
		err = errors.New("Can only transpose a square matrix in place")
		return
	}
	var i, j uint
	for i = 0; i < A.rows; i++ {
		for j = 0; j < i; j++ {
			tmp := A.Get(i, j)
			A.Set(i, j, A.Get(j, i))
			A.Set(j, i, tmp)
		}
	}
	return
}

func solveLower(A *DenseMatrix, b Matrix) *DenseMatrix {
	x := make([]float64, A.Cols())
	var i, j uint
	for i = 0; i < A.Rows(); i++ {
		x[i] = b.Get(i, 0)
		for j = 0; j < i; j++ {
			x[i] -= x[j] * A.Get(i, j)
		}
		//the diagonal defined to be ones
		//x[i] /= A.Get(i, i);
	}
	return MakeDenseMatrix(x, A.Cols(), 1)
}

func solveUpper(A *DenseMatrix, b Matrix) *DenseMatrix {
	x := make([]float64, A.Cols())
	for i := A.Rows() - 1; i >= 0; i-- {
		x[i] = b.Get(i, 0)
		for j := i + 1; j < A.Cols(); j++ {
			x[i] -= x[j] * A.Get(i, j)
		}
		x[i] /= A.Get(i, i)
	}
	return MakeDenseMatrix(x, A.Cols(), 1)
}

func (A *DenseMatrix) Solve(b MatrixRO) (*DenseMatrix, error) {
	Acopy := A.Copy()
	P := Acopy.LUInPlace()
	Pinv := P.Inverse()
	pb, err := Pinv.Times(b)

	if !(err == nil) {
		return nil, err
	}

	y := solveLower(Acopy, pb)
	x := solveUpper(Acopy, y)
	return x, nil
}

func (A *DenseMatrix) SolveDense(b *DenseMatrix) (*DenseMatrix, error) {
	return A.Solve(b)
}
