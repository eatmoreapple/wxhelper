package apiclient

import (
	"context"
	"encoding/base64"
	"encoding/json"
	. "github.com/eatmoreapple/wxhelper/internal/models"
	"io"
	"net/http"
)

type Client struct {
	transport *Transport
}

func (c *Client) GetUserInfo(ctx context.Context) (*Account, error) {
	resp, err := c.transport.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r Result[*Account]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if err = r.Err(); err != nil {
		return nil, err
	}
	return r.Data, nil
}

func (c *Client) CheckLogin(ctx context.Context) (bool, error) {
	resp, err := c.transport.CheckLogin(ctx)
	if err != nil {
		return false, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r Result[bool]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return false, err
	}
	if err = r.Err(); err != nil {
		return false, err
	}
	return r.Data, nil
}

func (c *Client) GetContactList(ctx context.Context) (Members, error) {
	resp, err := c.transport.GetContactList(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r Result[Members]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if err = r.Err(); err != nil {
		return nil, err
	}
	return r.Data, nil
}

func (c *Client) SendText(ctx context.Context, to, content string) error {
	resp, err := c.transport.SendText(ctx, to, content)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r Result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	return r.Err()
}

func (c *Client) SendImage(ctx context.Context, to string, img io.Reader) error {
	// to base64
	data, err := io.ReadAll(img)
	if err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(data)

	resp, err := c.transport.SendImage(ctx, to, encoded)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r Result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	return r.Err()
}

func (c *Client) SendFile(ctx context.Context, to string, file io.Reader) error {
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(data)

	resp, err := c.transport.SendFile(ctx, to, encoded)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r Result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	return r.Err()
}

func (c *Client) SyncMessage(ctx context.Context) ([]*Message, error) {
	resp, err := c.transport.SyncMessage(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r Result[[]*Message]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return r.Data, nil
}

func (c *Client) GetChatRoomDetail(ctx context.Context, chatRoomId string) (*ChatRoomInfo, error) {
	resp, err := c.transport.GetChatRoomDetail(ctx, chatRoomId)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r Result[ChatRoomInfo]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if err = r.Err(); err != nil {
		return nil, err
	}
	return &r.Data, nil
}

func (c *Client) GetMemberFromChatRoom(ctx context.Context, chatRoomId string) ([]*Profile, error) {
	resp, err := c.transport.GetMemberFromChatRoom(ctx, chatRoomId)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r Result[[]*Profile]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if err = r.Err(); err != nil {
		return nil, err
	}
	return r.Data, nil
}

func (c *Client) SendAtText(ctx context.Context, opt SendAtTextOption) error {
	resp, err := c.transport.SendAtText(ctx, opt)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r Result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	return r.Err()
}

func (c *Client) AddMemberIntoChatRoom(ctx context.Context, chatRoomID string, memberIDs []string) error {
	resp, err := c.transport.AddMemberIntoChatRoom(ctx, chatRoomID, memberIDs)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r Result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	return r.Err()
}

func New(apiServerURL string) *Client {
	return &Client{
		transport: &Transport{
			baseURL:    apiServerURL,
			httpClient: &http.Client{},
		},
	}
}
