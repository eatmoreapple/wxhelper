package apiclient

import "errors"

var ErrAuth = errors.New("auth error")

const (
	authErrCode = 2
)

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data,omitempty"`
}

func (r Result[T]) OK() bool {
	return r.Code == 0
}

func (r Result[T]) Err() error {
	if r.OK() {
		return nil
	}
	switch r.Code {
	case authErrCode:
		return ErrAuth
	default:
		return errors.New(r.Msg)
	}
}
