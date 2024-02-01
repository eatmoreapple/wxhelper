package apiclient

import (
	"context"
	"github.com/eatmoreapple/wxhelper/apiserver"
	"net/http"
	urlpkg "net/url"
)

type Transport struct {
	BaseURL string
}

// GetUserInfo GetUserInfo
func (c *Transport) GetUserInfo(ctx context.Context) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + apiserver.GetUserInfo)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req.WithContext(ctx))
}

// CheckLogin CheckLogin
func (c *Transport) CheckLogin(ctx context.Context) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + apiserver.CheckLogin)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req.WithContext(ctx))
}

func (c *Transport) GetContactList(ctx context.Context) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + apiserver.GetContactList)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req.WithContext(ctx))
}
