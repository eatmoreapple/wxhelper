package wxhelper

import (
	"context"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Client interface {
	// CheckLogin checks whether the user has logged in.
	CheckLogin(ctx context.Context) (bool, error)

	// GetUserInfo gets the user's information.
	GetUserInfo(ctx context.Context) (*Account, error)

	// SendText sends a text message to the user.
	SendText(ctx context.Context, to string, content string) error

	// SendImage sends an image message to the user.
	SendImage(ctx context.Context, to string, img io.Reader) error

	// GetContactList gets the contact list.
	GetContactList(context.Context) (Members, error)

	// HTTPHookSyncMsg hooks the sync message.
	HTTPHookSyncMsg(ctx context.Context, o HookSyncMsgOption) error

	// HTTPUnhookSyncMsg unhooks the sync message.
	HTTPUnhookSyncMsg(ctx context.Context) error
}

type client struct {
	transport *Transport
}

func (c *client) CheckLogin(ctx context.Context) (bool, error) {
	resp, err := c.transport.CheckLogin(ctx)
	if err != nil {
		return false, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[string]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return false, err
	}
	return r.OK(), r.Err()
}

func (c *client) GetUserInfo(ctx context.Context) (*Account, error) {
	resp, err := c.transport.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[Account]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r.Data, r.Err()
}

func (c *client) SendText(ctx context.Context, to string, content string) error {
	resp, err := c.transport.SendText(ctx, to, content)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	return r.Err()
}

func (c *client) GetContactList(ctx context.Context) (Members, error) {
	resp, err := c.transport.GetContactList(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[Members]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return r.Data, r.Err()
}

type HookSyncMsgOption struct {
	ServerURL *url.URL
	LocalURL  *url.URL
	Timeout   time.Duration
}

func (c *client) HTTPHookSyncMsg(ctx context.Context, o HookSyncMsgOption) error {
	opt := TransportHookSyncMsgOption{
		Port:       o.ServerURL.Port(),
		Ip:         o.ServerURL.Hostname(),
		Url:        o.LocalURL.String(),
		EnableHttp: true,
		Timeout:    strconv.Itoa(int(o.Timeout.Seconds()) * 100),
	}
	resp, err := c.transport.HookSyncMsg(ctx, opt)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {

		return err
	}
	return r.Err()
}

func (c *client) HTTPUnhookSyncMsg(ctx context.Context) error {
	resp, err := c.transport.UnhookSyncMsg(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	return r.Err()
}

func (c *client) SendImage(ctx context.Context, to string, img io.Reader) error {
	file, cb, err := readerToFile(img)
	if err != nil {
		return err
	}
	defer cb()
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	resp, err := c.transport.SendImage(ctx, to, stat.Name())
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	return r.Err()
}

func NewClient(transport *Transport) Client {
	return &client{transport: transport}
}

func DefaultClient() Client {
	transport := NewTransport(defaultServerURL.String())
	return NewClient(transport)
}

func readerToFile(reader io.Reader) (file *os.File, cb func(), err error) {
	var ok bool
	if file, ok = reader.(*os.File); ok {
		return file, func() {}, nil
	}
	file, err = os.CreateTemp("", "*")
	if err != nil {
		return nil, nil, err
	}
	cb = func() {
		_ = file.Close()
		_ = os.Remove(file.Name())
	}
	_, err = io.Copy(file, reader)
	if err != nil {
		cb()
		return nil, nil, err
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		cb()
		return nil, nil, err
	}
	return file, cb, nil
}
