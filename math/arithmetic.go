package math

//import "math"

/*
Tests the element-wise equality of the two matrices.
*/
func Equals(A, B MatrixRO) bool {
	if A.Rows() != B.Rows() || A.Cols() != B.Cols() {
		return false
	}
    var i, j uint
	for i = 0; i < A.Rows(); i++ {
		for j = 0; j < A.Cols(); j++ {
			if A.Get(i, j) != B.Get(i, j) {
				return false
			}
		}
	}
	return true
}
