package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/eatmoreapple/wxhelper/apiserver"
	"io"
	"mime/multipart"
	"net/http"
	urlpkg "net/url"
	"strconv"
)

type Transport struct {
	baseURL    string
	httpClient *http.Client
}

// GetUserInfo 获取用户信息
func (c *Transport) GetUserInfo(ctx context.Context) (*http.Response, error) {
	url, err := urlpkg.Parse(c.baseURL + apiserver.GetUserInfo)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}

// CheckLogin 检查是否登录
func (c *Transport) CheckLogin(ctx context.Context) (*http.Response, error) {
	url, err := urlpkg.Parse(c.baseURL + apiserver.CheckLogin)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req.WithContext(ctx))
}

// GetContactList 获取联系人列表
func (c *Transport) GetContactList(ctx context.Context) (*http.Response, error) {
	url, err := urlpkg.Parse(c.baseURL + apiserver.GetContactList)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req.WithContext(ctx))
}

// SendText 发送文本消息
func (c *Transport) SendText(ctx context.Context, to, content string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.baseURL + apiserver.SendText)
	if err != nil {
		return nil, err
	}
	var payload = apiserver.SendTextRequest{
		To:      to,
		Content: content,
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
	return c.httpClient.Do(req)
}

func (c *Transport) SendImage(ctx context.Context, to, filename string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.baseURL + apiserver.SendImage)
	if err != nil {
		return nil, err
	}
	var payload = apiserver.SendImageRequest{
		To:    to,
		Image: filename,
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
	return c.httpClient.Do(req)
}

func (c *Transport) SendFile(ctx context.Context, to, file string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.baseURL + apiserver.SendFile)
	if err != nil {
		return nil, err
	}
	var payload = apiserver.SendFileRequest{
		To:   to,
		File: file,
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
	return c.httpClient.Do(req)
}

// SyncMessage SyncMessage
func (c *Transport) SyncMessage(ctx context.Context) (*http.Response, error) {
	url, err := urlpkg.Parse(c.baseURL + apiserver.SyncMessage)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}

func (c *Transport) GetChatRoomDetail(ctx context.Context, chatRoomID string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.baseURL + apiserver.GetChatRoomDetail)
	if err != nil {
		return nil, err
	}
	var payload = apiserver.GetChatRoomInfoRequest{
		ChatRoomID: chatRoomID,
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
	return c.httpClient.Do(req)
}

func (c *Transport) GetMemberFromChatRoom(ctx context.Context, chatRoomID string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.baseURL + apiserver.GetMemberFromChatRoom)
	if err != nil {
		return nil, err
	}
	var payload = apiserver.GetMemberFromChatRoomRequest{
		ChatRoomID: chatRoomID,
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
	return c.httpClient.Do(req)
}

type SendAtTextOption struct {
	GroupID string
	Content string
	AtList  []string
}

func (c *Transport) SendAtText(ctx context.Context, option SendAtTextOption) (*http.Response, error) {
	url, err := urlpkg.Parse(c.baseURL + apiserver.SendAtText)
	if err != nil {
		return nil, err
	}
	var payload = apiserver.SendAtTextRequest{
		GroupID: option.GroupID,
		AtList:  option.AtList,
		Content: option.Content,
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
	return c.httpClient.Do(req)
}

func (c *Transport) AddMemberIntoChatRoom(ctx context.Context, chatRoomID string, memberIDs []string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.baseURL + apiserver.AddMemberToChatRoom)
	if err != nil {
		return nil, err
	}
	var payload = apiserver.AddMemberToChatRoomRequest{
		ChatRoomID: chatRoomID,
		MemberIds:  memberIDs,
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
	return c.httpClient.Do(req)
}

func (c *Transport) InviteMemberToChatRoom(ctx context.Context, chatRoomID string, memberIDs []string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.baseURL + apiserver.InviteMemberToChatRoom)
	if err != nil {
		return nil, err
	}
	var payload = apiserver.InviteMemberToChatRoomRequest{
		ChatRoomID: chatRoomID,
		MemberIds:  memberIDs,
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
	return c.httpClient.Do(req)
}

func (c *Transport) ForwardMsg(ctx context.Context, wxID, msgID string) (*http.Response, error) {
	url, err := urlpkg.Parse(c.baseURL + apiserver.ForwardMsg)
	if err != nil {
		return nil, err
	}
	var payload = apiserver.ForwardMsgRequest{
		WxID:  wxID,
		MsgID: msgID,
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
	return c.httpClient.Do(req)
}

func (c *Transport) UploadFile(ctx context.Context, request apiserver.UploadRequest) (*http.Response, error) {
	url, err := urlpkg.Parse(c.baseURL + apiserver.UploadFile)
	if err != nil {
		return nil, err
	}
	var buf = new(bytes.Buffer)

	writer := multipart.NewWriter(buf)

	part, err := writer.CreateFormFile("file", request.Filename)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(part, request.Content); err != nil {
		return nil, err
	}
	if err = writer.WriteField("filename", request.Filename); err != nil {
		return nil, err
	}
	if err = writer.WriteField("fileHash", request.FileHash); err != nil {
		return nil, err
	}
	if err = writer.WriteField("chunks", strconv.Itoa(request.Chunks)); err != nil {
		return nil, err
	}
	if err = writer.WriteField("chunk", strconv.Itoa(request.Chunk)); err != nil {
		return nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return c.httpClient.Do(req)
}
