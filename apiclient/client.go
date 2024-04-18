package apiclient

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/eatmoreapple/wxhelper/apiserver"
	. "github.com/eatmoreapple/wxhelper/internal/models"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
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
	path, err := c.UploadFile(ctx, "", img)
	if err != nil {
		return err
	}
	resp, err := c.transport.SendImage(ctx, to, path)
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
	path, err := c.UploadFile(ctx, "", file)
	if err != nil {
		return err
	}
	resp, err := c.transport.SendFile(ctx, to, path)
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

func (c *Client) InviteMemberToChatRoom(ctx context.Context, chatRoomID string, memberIDs []string) error {
	resp, err := c.transport.InviteMemberToChatRoom(ctx, chatRoomID, memberIDs)
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

func (c *Client) ForwardMsg(ctx context.Context, wxID, msgID string) error {
	resp, err := c.transport.ForwardMsg(ctx, wxID, msgID)
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

func (c *Client) UploadFile(ctx context.Context, filename string, reader io.Reader) (string, error) {
	if len(filename) == 0 {
		if f, ok := reader.(*os.File); ok {
			filename = f.Name()
		} else {
			filename = uuid.New().String()
		}
	}
	tmpFile, err := os.CreateTemp("", "*")
	if err != nil {
		return "", err
	}
	defer func() { _ = tmpFile.Close() }()

	h := sha256.New()

	if _, err = io.Copy(tmpFile, io.TeeReader(reader, h)); err != nil {
		return "", err
	}
	fileHash := hex.EncodeToString(h.Sum(nil))

	stat, err := tmpFile.Stat()
	if err != nil {
		return "", err
	}

	const chunkSize int64 = (1 << 20) / 2

	chunks := (stat.Size() + chunkSize - 1) / chunkSize

	// closure function to upload file
	upload := func(chunk int, reader io.Reader) (string, error) {
		resp, err := c.transport.UploadFile(ctx, apiserver.UploadRequest{
			Filename: filename,
			FileHash: fileHash,
			Chunks:   int(chunks),
			Chunk:    chunk,
			Content:  io.NopCloser(reader),
		})
		if err != nil {
			return "", err
		}
		defer func() { _ = resp.Body.Close() }()
		var r Result[string]
		if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return "", err
		}
		return r.Data, nil
	}

	var result string
	for i := int64(0); i < chunks; i++ {
		sectionReader := io.NewSectionReader(tmpFile, i*chunkSize, chunkSize)
		result, err = upload(int(i), sectionReader)
		if err != nil {
			return "", err
		}
	}
	return result, nil
}

func New(apiServerURL string) *Client {
	return &Client{
		transport: &Transport{
			baseURL:    apiServerURL,
			httpClient: &http.Client{},
		},
	}
}
