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
	url, err := urlpkg.Parse(c.BaseURL + apiserver.CheckLogin)
}
