package wxclient

import (
	"context"
	"encoding/json"
	"errors"
	. "github.com/eatmoreapple/wxhelper/internal/models"
	"io"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Client struct {
	transport *Transport
}

func (c *Client) CheckLogin(ctx context.Context) (bool, error) {
	resp, err := c.transport.CheckLogin(ctx)
	if err != nil {
		return false, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return false, err
	}
	return r.Code == 1, nil
}

func (c *Client) GetUserInfo(ctx context.Context) (*Account, error) {
	resp, err := c.transport.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[*Account]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if r.Code != 1 {
		return nil, errors.New("get user info failed")
	}
	return r.Data, nil
}

func (c *Client) SendText(ctx context.Context, to string, content string) error {
	resp, err := c.transport.SendText(ctx, to, content)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	if r.Code == 0 {
		return errors.New("send text failed")
	}
	return nil
}

func (c *Client) GetContactList(ctx context.Context) (Members, error) {
	resp, err := c.transport.GetContactList(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[Members]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return r.Data, nil
}

type HookSyncMsgOption struct {
	LocalURL *url.URL
	Timeout  time.Duration
}

func (c *Client) HTTPHookSyncMsg(ctx context.Context, o HookSyncMsgOption) error {
	opt := TransportHookSyncMsgOption{
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
	if r.Code != 0 {
		return errors.New("hook sync msg failed")
	}
	return nil
}

func (c *Client) HTTPUnhookSyncMsg(ctx context.Context) error {
	resp, err := c.transport.UnhookSyncMsg(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	if r.Code != 0 {
		return errors.New("unhook sync msg failed")
	}
	return nil
}

func (c *Client) SendImage(ctx context.Context, to string, img io.Reader) error {
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
	return nil
}

func New(transport *Transport) *Client {
	return &Client{transport: transport}
}

func Default() *Client {
	transport := NewTransport(os.Getenv("VIRTUAL_MACHINE_URL"))
	return New(transport)
}
