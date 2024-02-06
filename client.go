package wxhelper

import (
	"context"
	"errors"
	"github.com/eatmoreapple/wxhelper/apiclient"
	"github.com/eatmoreapple/wxhelper/pkg/structcopy"
	"io"
)

var ErrNotLogin = errors.New("user not login")

type Client struct {
	apiclient *apiclient.Client
}

func (c *Client) GetUserInfo(ctx context.Context) (*Account, error) {
	account, err := c.apiclient.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	if account.Account == "" {
		return nil, ErrNotLogin
	}
	return structcopy.Copy[*Account](account)
}

func (c *Client) GetContactList(ctx context.Context) (Members, error) {
	members, err := c.apiclient.GetContactList(ctx)
	if err != nil {
		return nil, err
	}
	return structcopy.CopySlice[*User](members)
}

func (c *Client) SendText(ctx context.Context, to, content string) error {
	return c.apiclient.SendText(ctx, to, content)
}

func (c *Client) CheckLogin(ctx context.Context) (bool, error) {
	return c.apiclient.CheckLogin(ctx)
}

func (c *Client) SendImage(ctx context.Context, to string, img io.Reader) error {
	return c.apiclient.SendImage(ctx, to, img)
}

func (c *Client) SendFile(ctx context.Context, to string, file io.Reader) error {
	return c.apiclient.SendFile(ctx, to, file)
}

func (c *Client) SyncMessage(ctx context.Context) ([]*Message, error) {
	message, err := c.apiclient.SyncMessage(ctx)
	if err != nil {
		return nil, err
	}
	return structcopy.CopySlice[*Message](message)
}

func (c *Client) GetChatRoomInfo(ctx context.Context, chatRoomID string) (*GroupInfo, error) {
	chatRoomInfo, err := c.apiclient.GetChatRoomDetail(ctx, chatRoomID)
	if err != nil {
		return nil, err
	}
	return structcopy.Copy[*GroupInfo](chatRoomInfo)
}

func (c *Client) GetChatRoomMembers(ctx context.Context, chatRoomID string) ([]*Profile, error) {
	members, err := c.apiclient.GetMemberFromChatRoom(ctx, chatRoomID)
	if err != nil {
		return nil, err
	}
	return structcopy.CopySlice[*Profile](members)
}
