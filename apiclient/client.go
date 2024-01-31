package apiclient

import (
	"context"
	. "github.com/eatmoreapple/wxhelper/apiserver/models"
)

type Client struct{}

func (c *Client) GetUserInfo(ctx context.Context) (*Account, error) {
	return nil, nil
}

func (c *Client) GetContactList(ctx context.Context) (Members, error) {
	return nil, nil
}
