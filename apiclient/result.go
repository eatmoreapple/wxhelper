package apiclient

import "errors"

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
	return errors.New(r.Msg)
}
