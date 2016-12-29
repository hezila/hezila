package math

//returns a copy of the row (not a slice)
func (A *DenseMatrix) RowCopy(i uint) []float64 {
	row := make([]float64, A.cols)
    var j uint
	for j = 0; j < A.cols; j++ {
		row[j] = A.Get(i, j)
	}
	return row
}

//returns a copy of the column (not a slice)
func (A *DenseMatrix) ColCopy(j uint) []float64 {
	col := make([]float64, A.rows)
    var i uint
	for i = 0; i < A.rows; i++ {
		col[i] = A.Get(i, j)
	}
	return col
}

//returns a copy of the diagonal (not a slice)
func (A *DenseMatrix) DiagonalCopy() []float64 {
	span := A.rows
	if A.cols < span {
		span = A.cols
	}
	diag := make([]float64, span)
    var i uint
	for i = 0; i < span; i++ {
		diag[i] = A.Get(i, i)
	}
	return diag
}

func (A *DenseMatrix) BufferRow(i uint, buf []float64) {
    var j uint
	for j = 0; j < A.cols; j++ {
		buf[j] = A.Get(i, j)
	}
}

func (A *DenseMatrix) BufferCol(j uint, buf []float64) {
    var i uint
	for i = 0; i < A.rows; i++ {
		buf[i] = A.Get(i, j)
	}
}

func (A *DenseMatrix) BufferDiagonal(buf []float64) {
    var i uint
	for i = 0; i < A.rows && i < A.cols; i++ {
		buf[i] = A.Get(i, i)
	}
}

func (A *DenseMatrix) FillRow(i uint, buf []float64) {
    var j uint
	for j = 0; j < A.cols; j++ {
		A.Set(i, j, buf[j])
	}
}

func (A *DenseMatrix) FillCol(j uint, buf []float64) {
    var i uint
	for i = 0; i < A.rows; i++ {
		A.Set(i, j, buf[i])
	}
}

func (A *DenseMatrix) FillDiagonal(buf []float64) {
    var i uint
	for i = 0; i < A.rows && i < A.cols; i++ {
		A.Set(i, i, buf[i])
	}
}
