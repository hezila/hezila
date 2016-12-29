package math

import "math"

/*
Swap two rows in this PivotMatrix.
*/
func (P *PivotMatrix) SwapRows(r1, r2 uint) error {
	//	tmp := P.pivots[r1];
	//	P.pivots[r1] = P.pivots[r2];
	//	P.pivots[r2] = tmp;
	P.pivots[r1], P.pivots[r2] = P.pivots[r2], P.pivots[r1]
	P.pivotSign *= -1

	return nil
}

func (P *PivotMatrix) Symmetric() bool {
    var i uint
	for i = 0; i < P.rows; i++ {
		if P.pivots[P.pivots[i]] != i {
			return false
		}
	}
	return true
}

func (A *PivotMatrix) Inverse() *PivotMatrix { return A.Transpose() }

func (P *PivotMatrix) Transpose() *PivotMatrix {
	newPivots := make([]int, P.rows)
    var i uint
	for i = 0; i < P.rows; i++ {
		newPivots[P.pivots[i]] = i
	}
	return MakePivotMatrix(newPivots, P.pivotSign)
}

func (P *PivotMatrix) Det() float64 { return P.pivotSign }

func (P *PivotMatrix) Trace() (r float64) {
    var i uint
	for i = 0; i < len(P.pivots); i++ {
		if P.pivots[i] == i {
			r += 1
		}
	}
	return
}

/*
Returns x such that Px=b.
*/
func (P *PivotMatrix) Solve(b MatrixRO) (Matrix, error) {
	return P.Transpose().Times(b) //error comes from times
}

func (A *PivotMatrix) OneNorm() float64      { return float64(A.rows) }
func (A *PivotMatrix) TwoNorm() float64      { return math.Sqrt(float64(A.rows)) }
func (A *PivotMatrix) InfinityNorm() float64 { return 1 }

