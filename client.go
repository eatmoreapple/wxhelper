package wxhelper

import (
	"context"
	"github.com/eatmoreapple/wxhelper/apiclient"
	"github.com/eatmoreapple/wxhelper/pkg/structcopy"
	"io"
)

type Client struct {
	apiclient *apiclient.Client
}

func (c *Client) GetUserInfo(ctx context.Context) (*Account, error) {
	account, err := c.apiclient.GetUserInfo(ctx)
	if err != nil {
		return nil, err
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

func (c *Client) SyncMessage(ctx context.Context) ([]*Message, error) {
	message, err := c.apiclient.SyncMessage(ctx)
	if err != nil {
		return nil, err
	}
	return structcopy.CopySlice[*Message](message)
}
