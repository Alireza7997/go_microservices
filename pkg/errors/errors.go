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

// StatusError is an error that carries an HTTP status code and a
// client-safe response message.
type StatusError interface {
	error
	StatusCode() int64
	Message() string
}

func (e *theError) StatusCode() int64 {
	return e.Code
}

// Message returns the client-safe response message.
func (e *theError) Message() string {
	return e.ResponseMessage
}

func New(code int64, errMessage, responseMessage string, err ...error) error {
	var wrapped error
	if len(err) > 0 {
		wrapped = err[0]
	}

	return &theError{
		Code:            code,
		ErrMessage:      errMessage,
		ResponseMessage: responseMessage,
		Err:             wrapped,
	}
}
