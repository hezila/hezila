package math

import "math"

/*
Overwrites A with [L\U] and returns P, st PLU=A. L is considered to
have 1s in the diagonal.
*/
func (A *DenseMatrix) LUInPlace() (P *PivotMatrix) {
	m := A.Rows()
	n := A.Cols()
	LUcolj := make([]float64, m)
	LUrowi := make([]float64, n)
	piv := make([]uint, m)
	var i, j, k uint
	for i = 0; i < m; i++ {
		piv[i] = i
	}
	pivsign := float64(1.0)

	for j = 0; j < n; j++ {
		A.BufferCol(j, LUcolj)
		for i = 0; i < m; i++ {
			A.BufferRow(i, LUrowi)
			kmax := i
			if j < i {
				kmax = j
			}
			s := float64(0)
			for k = 0; k < kmax; k++ {
				s += LUrowi[k] * LUcolj[k]
			}
			LUcolj[i] -= s
			LUrowi[j] = LUcolj[i]
			A.Set(i, j, LUrowi[j])
		}

		p := j
		for i := j + 1; i < m; i++ {
			if math.Abs(LUcolj[i]) > math.Abs(LUcolj[p]) {
				p = i
			}
		}
		if p != j {
			A.SwapRows(p, j)
			k := piv[p]
			piv[p] = piv[j]
			piv[j] = k
			pivsign = -pivsign
		}

		if j < m && A.Get(j, j) != 0 {
			for i := j + 1; i < m; i++ {
				A.Set(i, j, A.Get(i, j)/A.Get(j, j))
			}
		}
	}

	P = MakePivotMatrix(piv, pivsign)

	return
}
