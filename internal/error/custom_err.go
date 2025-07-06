package cerr

import "strings"

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func New(code int, msg ...string) *Error {
	var message string
	if len(msg) <= 0 {
		message = GetMessage(code)
	} else {
		message = strings.Join(msg, " ")
	}
	return &Error{
		Code:    code,
		Message: message,
	}
}

func NewSuccess(msg ...string) *Error {
	return New(SUCCESS, msg...)
}

func NewSysError(msg ...string) *Error {
	return New(SYSTEM_ERROR, msg...)
}

func NewParamError(msg ...string) *Error {
	return New(BAD_REQUEST, msg...)
}
