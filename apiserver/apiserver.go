package apiserver

import (
	"context"
	"errors"
	"github.com/eatmoreapple/ginx"
	"github.com/eatmoreapple/wxhelper/apiserver/internal/msgbuffer"
	. "github.com/eatmoreapple/wxhelper/apiserver/models"
	"github.com/eatmoreapple/wxhelper/internal/wxclient"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"time"
)

// APIServer 用来屏蔽微信的接口

type APIServer struct {
	client      *wxclient.Client
	msgBuffer   msgbuffer.MessageBuffer
	msgListener MessageListener
	engine      *gin.Engine
}

// CheckLogin 检查是否登录
func (a *APIServer) CheckLogin(ctx context.Context, _ struct{}) (*Result[bool], error) {
	ok, err := a.client.CheckLogin(ctx)
	if err != nil {
		return nil, err
	}
	return OK(ok), nil
}

// GetUserInfo 获取用户信息
func (a *APIServer) GetUserInfo(ctx context.Context, _ struct{}) (*Result[*Account], error) {
	account, err := a.client.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	// FIXME
	entity := (*Account)(account)
	return OK(entity), nil
}

type SendTextRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

// SendText 发送文本消息
func (a *APIServer) SendText(ctx context.Context, req SendTextRequest) (*Result[any], error) {
	err := a.client.SendText(ctx, req.To, req.Content)
	if err != nil {
		return nil, err
	}
	return OK[any](nil), nil
}

// GetContactList 获取联系人列表
func (a *APIServer) GetContactList(ctx context.Context, _ struct{}) (*Result[Members], error) {
	members, err := a.client.GetContactList(ctx)
	if err != nil {
		return nil, err
	}
	var entity = make(Members, 0, len(members))
	for _, member := range members {
		// FIXME
		entity = append(entity, (*User)(member))
	}
	return OK(entity), nil
}

// SyncMessage 同步消息
func (a *APIServer) SyncMessage(ctx context.Context, _ struct{}) (*Result[[]*Message], error) {
	message, err := a.msgBuffer.Get(ctx, time.Second*25)
	if errors.Is(err, msgbuffer.ErrTimeout) {
		messages := make([]*Message, 0)
		return OK(messages), nil
	}
	if err != nil {
		return nil, err
	}
	return OK([]*Message{message}), nil
}

func (a *APIServer) startListen() error {
	listenURL, _ := url.Parse("http://localhost:9999")
	go func() { _ = http.ListenAndServe(":9999", a.msgListener) }()
	time.Sleep(time.Second * 1)
	return a.client.HTTPHookSyncMsg(context.Background(), wxclient.HookSyncMsgOption{
		LocalURL: listenURL,
		Timeout:  time.Second * 30,
	})
}

func (a *APIServer) Run() error {
	router := ginx.NewRouter(a.engine)
	registerAPIServer(router, a)
	if err := a.startListen(); err != nil {
		return err
	}
	return a.engine.Run(":19089")
}

func New(client *wxclient.Client, msgBuffer msgbuffer.MessageBuffer) *APIServer {
	srv := &APIServer{
		client:      client,
		msgBuffer:   msgBuffer,
		msgListener: NewMessageListener(msgBuffer),
		engine:      gin.New(),
	}
	return srv
}

func Default() *APIServer {
	return New(wxclient.Default(), msgbuffer.Default())
}
