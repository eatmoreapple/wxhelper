package wxclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eatmoreapple/env"
	. "github.com/eatmoreapple/wxhelper/internal/models"
	"io"
	"net/url"
	"strconv"
	"strings"
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

func (c *Client) HTTPHookSyncMsg(ctx context.Context, url *url.URL, timeout time.Duration) error {
	opt := TransportHookSyncMsgOption{
		Url:        url.String(),
		EnableHttp: 1,
		Timeout:    strconv.Itoa(int(timeout / time.Second)),
		Ip:         url.Hostname(),
		Port:       url.Port(),
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

func (c *Client) HookSyncMsg(ctx context.Context, ip string, port int) error {
	opt := TransportHookSyncMsgOption{
		EnableHttp: 0,
		Ip:         ip,
		Port:       strconv.Itoa(port),
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

func (c *Client) UnhookSyncMsg(ctx context.Context) error {
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
	filename, err := saveToLocal(img)
	if err != nil {
		return err
	}
	resp, err := c.transport.SendImage(ctx, to, filename)
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

func (c *Client) SendFile(ctx context.Context, to string, img io.Reader) error {
	filename, err := saveToLocal(img)
	if err != nil {
		return err
	}
	resp, err := c.transport.SendFile(ctx, to, filename)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	if r.Code == 0 {
		return fmt.Errorf("send file failed with code %d", r.Code)
	}
	return nil
}

func (c *Client) GetChatRoomDetail(ctx context.Context, chatRoomId string) (*ChatRoomInfo, error) {
	resp, err := c.transport.GetChatRoomDetail(ctx, chatRoomId)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[ChatRoomInfo]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if r.Code != 1 {
		return nil, errors.New("get chat room detail failed")
	}
	return &r.Data, nil
}

func (c *Client) GetMemberFromChatRoom(ctx context.Context, chatRoomId string) (*GroupMember, error) {
	resp, err := c.transport.GetMemberFromChatRoom(ctx, chatRoomId)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[GroupMember]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if r.Code != 1 {
		return nil, errors.New("get chat room member failed")
	}
	return &r.Data, nil
}

func (c *Client) GetContactProfile(ctx context.Context, wxid string) (*Profile, error) {
	resp, err := c.transport.GetContactProfile(ctx, wxid)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[Profile]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if r.Code < 0 {
		return nil, errors.New("get contact profile failed")
	}
	return &r.Data, nil
}

type SendAtTextOption struct {
	WxIds      []string
	ChatRoomID string
	Content    string
}

func (c *Client) SendAtText(ctx context.Context, opt SendAtTextOption) error {
	resp, err := c.transport.SendAtText(ctx, sendAtTextOption{
		WxIds:      strings.Join(opt.WxIds, ","),
		ChatRoomId: opt.ChatRoomID,
		Msg:        opt.Content,
	})
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	if r.Code < 0 {
		return errors.New("send at text failed")
	}
	return nil
}

func (c *Client) AddMemberIntoChatRoom(ctx context.Context, chatRoomID string, memberIDs []string) error {
	resp, err := c.transport.AddMemberIntoChatRoom(ctx, chatRoomID, strings.Join(memberIDs, ","))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	if r.Code != 1 {
		return errors.New("add member into chat room failed")
	}
	return nil
}

func New(transport *Transport) *Client {
	return &Client{transport: transport}
}

func Default() *Client {
	transport := NewTransport(env.Name("VIRTUAL_MACHINE_URL").String())
	return New(transport)
}
