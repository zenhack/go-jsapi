package jsapi

import "syscall/js"

type Error Value

func (e Error) Error() string {
	return e.Value.Get("message").String()
}

func NewError(s string) Error {
	return Error{Value: js.Global().Get("Error").New(s)}
}

func WrapError(err error) Error {
	e, ok := err.(Error)
	if ok {
		return e
	}
	return NewError(err.Error())
}
