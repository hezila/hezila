package utils

// CommonError common error message in the format of json
type CommonError struct {
	ErrorCode int64  `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
}
