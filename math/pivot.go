package math

import (
	log "github.com/golang/glog"
)

/*
A space-optimized structure for pivot matrices, ie a matrix with
exactly one 1 in each row and each column.
*/
type PivotMatrix struct {
	matrix
	pivots    []uint
	pivotSign float64
}

// TODO: to implements
func (P *PivotMatrix) Arrays() [][]float64 {
	log.Warning("The Arrays() function for pivot matrix has not been implemented!")
	return nil
}

// TODO: to implements
func (P *PivotMatrix) Array() []float64 {
	log.Warning("The Array() function for pivot matrix has not been implemented!")
	return nil
}

func (P *PivotMatrix) Get(i, j uint) (v float64) {
	if i < 0 {
		i = P.rows + i
		if i < 0 {
			log.Fatal("index out of bound!")
		}
	}

	if j < 0 {
		j = P.cols + j
		if j < 0 {
			log.Fatal("index out of bound!")
		}
	}
	
	if P.pivots[j] == i {
		v = 1
	}
	return
}

/*
Convert this PivotMatrix into a DenseMatrix.
*/
func (P *PivotMatrix) DenseMatrix() *DenseMatrix {
	A := Zeros(P.rows, P.cols)
	for j := uint(0); j < P.rows; j++ {
		A.Set(P.pivots[j], j, 1)
	}
	return A
}

/*
Convert this PivotMatrix into a SparseMatrix.
*/
func (P *PivotMatrix) SparseMatrix() *SparseMatrix {
	A := ZerosSparse(P.rows, P.cols)

	for j := uint(0); j < P.rows; j++ {
		A.Set(P.pivots[j], j, 1)
	}
	return A
}

/*
Make a copy of this PivotMatrix.
*/
func (P *PivotMatrix) Copy() *PivotMatrix { return MakePivotMatrix(P.pivots, P.pivotSign) }

func MakePivotMatrix(pivots []uint, pivotSign float64) *PivotMatrix {
	n := uint(len(pivots))
	P := new(PivotMatrix)
	P.rows = n
	P.cols = n
	P.pivots = pivots
	P.pivotSign = pivotSign
	return P
}

func (A *PivotMatrix) String() string { return String(A) }
