package wxhelper

import "errors"

type result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data,omitempty"`
}

func (r result[T]) OK() bool {
	return r.Code != 0
}

func (r result[T]) Err() error {
	if r.OK() {
		return nil
	}
	return errors.New(r.Msg)
}
