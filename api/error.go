package api

import "fmt"

// HTTPError is an error with a message
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func httpError(code int, fmtString string, args ...interface{}) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: fmt.Sprintf(fmtString, args...),
	}
}
