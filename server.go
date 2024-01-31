package wxhelper

import (
	"context"
	"net/http"
	"net/url"
)

type Server interface {
	Register(ctx context.Context, node RegisterNode) error
}

type RegisterNode interface {
	URL() *url.URL
}

type MessageReceiver interface {
	ReceiveMessage(ctx context.Context) (<-chan *Message, error)
	RegisterNode
}

type HTTPRegisterNode struct {
	addr *url.URL
}

func (h *HTTPRegisterNode) URL() *url.URL {
	return h.addr
}

func (h *HTTPRegisterNode) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (h *HTTPRegisterNode) ReceiveMessage(ctx context.Context) (<-chan *Message, error) {
	panic("implement me")
}
