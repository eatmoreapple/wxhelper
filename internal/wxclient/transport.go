package wxclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	urlpkg "net/url"
)

type TransportHookSyncMsgOption struct {
	Port       string `json:"port"`
	Ip         string `json:"ip"`
	Url        string `json:"url"`
	Timeout    string `json:"timeout"`
	EnableHttp int    `json:"enableHttp"`
}

type Transport struct {
	BaseURL string
}

func (c *Transport) CheckLogin(ctx context.Context) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + "/api/checkLogin")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func (c *Transport) GetUserInfo(ctx context.Context) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + "/api/userInfo")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func (c *Transport) SendText(ctx context.Context, to string, content string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + "/api/sendTextMsg")
	if err != nil {
		return nil, err
	}
	var payload = map[string]string{
		"wxid": to,
		"msg":  content,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func (c *Transport) ForwardMessage(ctx context.Context, to, msgID string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + "/api/forwardMessage")
	if err != nil {
		return nil, err
	}
	var payload = map[string]string{
		"wxid":  to,
		"msgid": msgID,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func (c *Transport) SendImage(ctx context.Context, to, imagePath string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + "/api/sendImage")
	if err != nil {
		return nil, err
	}
	var payload = map[string]string{
		"wxid":      to,
		"imagePath": imagePath,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func (c *Transport) SendFile(ctx context.Context, to, filePath string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + "/api/sendFileMsg")
	if err != nil {
		return nil, err
	}
	var payload = map[string]string{
		"wxid":     to,
		"filePath": filePath,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func (c *Transport) GetContactList(ctx context.Context) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + "/api/getContactList")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func (c *Transport) HookSyncMsg(ctx context.Context, opt TransportHookSyncMsgOption) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + "/api/hookSyncMsg")
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(opt)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func (c *Transport) UnhookSyncMsg(ctx context.Context) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + "/api/unhookSyncMsg")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func NewTransport(baseURL string) *Transport {
	return &Transport{BaseURL: baseURL}
}
