package errors

import "fmt"

type theError struct {
	Code            int64
	ErrMessage      string
	ResponseMessage string
	Err             error
}

func (e *theError) Error() string {
	return fmt.Sprintf("Code: %d, ErrMessage: %s, ResponseMessage: %s, Err: %v", e.Code, e.ErrMessage, e.ResponseMessage, e.Err)
}

func New(code int64, errMessage, responseMessage string, err ...error) error {
	return &theError{
		Code:            code,
		ErrMessage:      errMessage,
		ResponseMessage: responseMessage,
		Err:             err[0],
	}
}
