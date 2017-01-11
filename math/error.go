package math

import (
	"fmt"
)

const (
	// The matrix returned was nil
	errorNilMatrix = iota
	// The element of indexed does not exist in (sparse) matrix
	errorNilElement
	// The dimensions of the inputs do not make sense
	errorDimensionMismatch
	// The indices provided are out of bounds
	errorIllegalIndex
	// The matrix provided has a singularity
	exceptionSingular
	// The matrix provided is not positive semi-definite.
	exceptionNotSPD
)

type error_ int

func (e error_) Error() string {
	switch e {
	case errorNilMatrix:
		return "Matrix is nil"
	case errorNilElement:
		return "Element indexed does not exist"
	case errorDimensionMismatch:
		return "Input dimensions do not match"
	case errorIllegalIndex:
		return "Index out of bounds"
	case exceptionSingular:
		return "Matrix is singular"
	case exceptionNotSPD:
		return "Matrix is not positive semidefinite"
	}
	return fmt.Sprintf("Unknown error code %d", e)
}

func (e error_) String() string {
	return e.Error()
}

var (
	//The matrix returned was nil.
	ErrorNilMatrix error_ = error_(errorNilMatrix)
	//The element identified in the matrix was nil (does not exist)
	ErrorNilElement error_ = error_(errorNilElement)
	//The dimensions of the inputs do not make sense for this operation.
	ErrorDimensionMismatch error_ = error_(errorDimensionMismatch)
	//The indices provided are out of bounds.
	ErrorIllegalIndex error_ = error_(errorIllegalIndex)
	//The matrix provided has a singularity.
	ExceptionSingular error_ = error_(exceptionSingular)
	//The matrix provided is not positive semi-definite.
	ExceptionNotSPD error_ = error_(exceptionNotSPD)
)
