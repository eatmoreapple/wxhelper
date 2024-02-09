package apiserver

type resultCode int

const (
	resultCodeOk resultCode = iota
	resultCodeErr
	resultCodeAuthErr
)

type Result[T any] struct {
	Code resultCode `json:"code"`
	Msg  string     `json:"msg"`
	Data T          `json:"data,omitempty"`
}

func OK[T any](data T) *Result[T] {
	return &Result[T]{
		Code: resultCodeOk,
		Data: data,
	}
}

func Err[T any](msg string) *Result[T] {
	return &Result[T]{
		Code: resultCodeErr,
		Msg:  msg,
	}
}
