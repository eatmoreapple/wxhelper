package wxhelper

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type MessageRetriever interface {
	RetrieveMessage() (<-chan *Message, error)
	Stop() error
}

type httpMessageRetriever struct {
	ctx         context.Context
	cancel      context.CancelFunc
	once        sync.Once
	messageChan chan *Message
	client      Client
	option      httpMessageRetrieverOption
}

type httpMessageRetrieverOption struct {
	messageChanSize uint
	serverURL       *url.URL
	localURL        *url.URL
	timeout         time.Duration
}

type HttpMessageRetrieverOptionFunc func(*httpMessageRetrieverOption)

func WithMessageChanSize(size uint) HttpMessageRetrieverOptionFunc {
	return func(o *httpMessageRetrieverOption) { o.messageChanSize = size }
}

func WithServerURL(u *url.URL) HttpMessageRetrieverOptionFunc {
	return func(o *httpMessageRetrieverOption) { o.serverURL = u }
}

func WithLocalURL(u *url.URL) HttpMessageRetrieverOptionFunc {
	return func(o *httpMessageRetrieverOption) { o.localURL = u }
}

func WithTimeout(t time.Duration) HttpMessageRetrieverOptionFunc {
	return func(o *httpMessageRetrieverOption) { o.timeout = t }
}

var (
	// defaultHttpMessageRetrieverOptions is the default options for httpMessageRetriever.
	defaultHttpMessageRetrieverOptions = []HttpMessageRetrieverOptionFunc{
		WithMessageChanSize(0),
		WithServerURL(defaultServerURL),
		WithLocalURL(defaultLocalURL),
		WithTimeout(time.Second * 30),
	}
)

func (r *httpMessageRetriever) RetrieveMessage() (<-chan *Message, error) {
	var err error
	r.once.Do(func() {
		r.messageChan = make(chan *Message, r.option.messageChanSize)
		srv := http.Server{Addr: ":" + r.option.localURL.Port(), Handler: r}
		defer func() { go func() { _ = srv.ListenAndServe() }() }()
		cancel := r.cancel
		r.cancel = func() {
			_ = r.client.HTTPUnhookSyncMsg(r.ctx)
			_ = srv.Shutdown(r.ctx)
			cancel()
		}
		_ = r.client.HTTPUnhookSyncMsg(r.ctx)
		opt := HookSyncMsgOption{ServerURL: r.option.serverURL, LocalURL: r.option.localURL, Timeout: time.Second * 30}
		err = r.client.HTTPHookSyncMsg(r.ctx, opt)
	})
	if err != nil {
		r.cancel()
		return nil, err
	}
	return r.messageChan, nil
}

func (r *httpMessageRetriever) Stop() error {
	r.cancel()
	close(r.messageChan)
	return nil
}

func (r *httpMessageRetriever) ServeHTTP(_ http.ResponseWriter, req *http.Request) {
	var msg Message
	if err := json.NewDecoder(req.Body).Decode(&msg); err != nil {
		return
	}
	select {
	case <-r.ctx.Done():
		return
	case r.messageChan <- &msg:
	}
}

var (
	defaultServerURL = &url.URL{Scheme: "http", Host: "localhost:19088"}
	defaultLocalURL  = &url.URL{Scheme: "http", Host: "172.24.176.1:19089"}
)

func NewHttpMessageRetriever(ctx context.Context, client Client, opts ...HttpMessageRetrieverOptionFunc) MessageRetriever {
	ctx, cancel := context.WithCancel(ctx)
	r := &httpMessageRetriever{
		ctx:    ctx,
		cancel: cancel,
		client: client,
	}
	opts = append(defaultHttpMessageRetrieverOptions, opts...)
	for _, optFunc := range opts {
		optFunc(&r.option)
	}
	return r
}
