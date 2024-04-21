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

type sendAtTextOption struct {
	ChatRoomId string `json:"chatRoomId"`
	WxIds      string `json:"wxids"`
	Msg        string `json:"msg"`
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
	url, err := urlpkg.Parse(c.BaseURL + "/api/sendImagesMsg")
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

func (c *Transport) GetChatRoomDetail(ctx context.Context, chatRoomId string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + "/api/getChatRoomDetailInfo")
	if err != nil {
		return nil, err
	}
	var payload = map[string]string{
		"chatRoomId": chatRoomId,
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

func (c *Transport) ModifyNickname(ctx context.Context, chatRoomId, wxId, nickname string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + "/api/modifyNickname")
	if err != nil {
		return nil, err
	}
	var payload = map[string]string{
		"chatRoomId": chatRoomId,
		"wxid":       wxId,
		"nickName":   nickname,
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

func (c *Transport) DelMemberFromChatRoom(ctx context.Context, chatRoomId string, memberIds ...string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + "/api/delMemberFromChatRoom")
	if err != nil {
		return nil, err
	}
	var payload = map[string]interface{}{
		"chatRoomId": chatRoomId,
		"memberIds":  memberIds,
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

func (c *Transport) GetMemberFromChatRoom(ctx context.Context, chatRoomId string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + "/api/getMemberFromChatRoom")
	if err != nil {
		return nil, err
	}
	var payload = map[string]string{
		"chatRoomId": chatRoomId,
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

func (c *Transport) GetContactProfile(ctx context.Context, wxid string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + "/api/getContactProfile")
	if err != nil {
		return nil, err
	}
	var payload = map[string]string{
		"wxid": wxid,
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

func (c *Transport) SendAtText(ctx context.Context, option sendAtTextOption) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + "/api/sendAtText")
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(option)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func (c *Transport) AddMemberIntoChatRoom(ctx context.Context, chatRoomId string, memberIds string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + "/api/addMemberToChatRoom")
	if err != nil {
		return nil, err
	}
	var payload = map[string]string{
		"chatRoomId": chatRoomId,
		"memberIds":  memberIds,
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

func (c *Transport) InviteMemberToChatRoom(ctx context.Context, chatRoomId string, memberIds string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + "/api/InviteMemberToChatRoom")
	if err != nil {
		return nil, err
	}
	var payload = map[string]string{
		"chatRoomId": chatRoomId,
		"memberIds":  memberIds,
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

func (c *Transport) ForwardMsg(ctx context.Context, msgID, wxID string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.BaseURL + "/api/forwardMsg")
	if err != nil {
		return nil, err
	}
	var payload = map[string]string{
		"wxid":  wxID,
		"msgId": msgID,
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

func NewTransport(baseURL string) *Transport {
	return &Transport{BaseURL: baseURL}
}
