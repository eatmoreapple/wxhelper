package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/eatmoreapple/wxhelper/apiserver"
	"net/http"
	urlpkg "net/url"
)

type Transport struct {
	BaseURL string
}

// GetUserInfo 获取用户信息
func (c *Transport) GetUserInfo(ctx context.Context) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + apiserver.GetUserInfo)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

// CheckLogin 检查是否登录
func (c *Transport) CheckLogin(ctx context.Context) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + apiserver.CheckLogin)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url.String(), nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req.WithContext(ctx))
}

// GetContactList 获取联系人列表
func (c *Transport) GetContactList(ctx context.Context) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + apiserver.GetContactList)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url.String(), nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req.WithContext(ctx))
}

// SendText 发送文本消息
func (c *Transport) SendText(ctx context.Context, to, content string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + apiserver.SendText)
	if err != nil {
		return nil, err
	}
	var payload = map[string]string{
		"to":      to,
		"content": content,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}

func (c *Transport) SendImage(ctx context.Context, to, imgData string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + apiserver.SendImage)
	if err != nil {
		return nil, err
	}
	var payload = map[string]string{
		"to":    to,
		"image": imgData,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req.WithContext(ctx))
}

// SyncMessage SyncMessage
func (c *Transport) SyncMessage(ctx context.Context) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + apiserver.SyncMessage)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url.String(), nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req.WithContext(ctx))
}
