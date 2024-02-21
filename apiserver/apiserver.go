package apiserver

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/eatmoreapple/env"
	"github.com/eatmoreapple/ginx"
	"github.com/eatmoreapple/wxhelper/apiserver/internal/msgbuffer"
	. "github.com/eatmoreapple/wxhelper/internal/models"
	"github.com/eatmoreapple/wxhelper/internal/wxclient"
	"github.com/eatmoreapple/wxhelper/pkg/netutil"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var ErrLogout = errors.New("logout")

type SendImageRequest struct {
	To          string `json:"to"`
	Image       string `json:"image"`
	imageReader io.Reader
}

type SendTextRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

func (a *SendImageRequest) FromContext(ctx *gin.Context) error {
	if err := json.NewDecoder(ctx.Request.Body).Decode(a); err != nil {
		return err
	}
	data, err := base64.StdEncoding.DecodeString(a.Image)
	if err != nil {
		return err
	}
	a.imageReader = bytes.NewReader(data)
	return nil
}

type SendFileRequest struct {
	To         string `json:"to"`
	Image      string `json:"file"`
	fileReader io.Reader
}

func (a *SendFileRequest) FromContext(ctx *gin.Context) error {
	if err := json.NewDecoder(ctx.Request.Body).Decode(a); err != nil {
		return err
	}
	data, err := base64.StdEncoding.DecodeString(a.Image)
	if err != nil {
		return err
	}
	a.fileReader = bytes.NewReader(data)
	return nil
}

type GetChatRoomInfoRequest struct {
	ChatRoomID string `json:"chatRoomId"`
}

type GetMemberFromChatRoomRequest struct {
	ChatRoomID string `json:"chatRoomId"`
}

type SendAtTextRequest struct {
	GroupID string   `json:"groupId"`
	AtList  []string `json:"atList"`
	Content string   `json:"content"`
}

type AddMemberToChatRoomRequest struct {
	ChatRoomID string   `json:"chatRoomId"`
	MemberIds  []string `json:"memberIds"`
}

type InviteMemberToChatRoomRequest struct {
	ChatRoomID string   `json:"chatRoomId"`
	MemberIds  []string `json:"memberIds"`
}

type ForwardMsgRequest struct {
	WxID  string `json:"wxid"`
	MsgID string `json:"msgId"`
}

// APIServer 用来屏蔽微信的接口
type APIServer struct {
	client    *wxclient.Client
	msgBuffer msgbuffer.MessageBuffer
	status    int32
	ctx       context.Context
	stop      context.CancelCauseFunc
	checker   Checker
}

func (a *APIServer) IsLogin() bool {
	return atomic.LoadInt32(&a.status) == 1
}

func (a *APIServer) logout() {
	atomic.StoreInt32(&a.status, 0)
	a.stop(ErrLogout)
}

func (a *APIServer) login() {
	atomic.StoreInt32(&a.status, 1)
}

func (a *APIServer) Ping(_ context.Context, _ ginx.Empty) (string, error) {
	return "pong", nil
}

// CheckLogin 检查是否登录
func (a *APIServer) CheckLogin(ctx context.Context, _ ginx.Empty) (*Result[bool], error) {
	ok, err := a.client.CheckLogin(ctx)
	if err != nil {
		return nil, err
	}
	return OK(ok), nil
}

// GetUserInfo 获取用户信息
func (a *APIServer) GetUserInfo(ctx context.Context, _ ginx.Empty) (*Result[*Account], error) {
	account, err := a.client.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return OK(account), nil
}

// SendText 发送文本消息
func (a *APIServer) SendText(ctx context.Context, req SendTextRequest) (*Result[any], error) {
	err := a.client.SendText(ctx, req.To, req.Content)
	if err != nil {
		return nil, err
	}
	return OK[any](nil), nil
}

func (a *APIServer) SendImage(ctx context.Context, req SendImageRequest) (*Result[any], error) {
	err := a.client.SendImage(ctx, req.To, req.imageReader)
	if err != nil {
		return nil, err
	}
	return OK[any](nil), nil
}

func (a *APIServer) SendFile(ctx context.Context, req SendFileRequest) (*Result[any], error) {
	err := a.client.SendFile(ctx, req.To, req.fileReader)
	if err != nil {
		return nil, err
	}
	return OK[any](nil), nil
}

// GetContactList 获取联系人列表
func (a *APIServer) GetContactList(ctx context.Context, _ ginx.Empty) (*Result[Members], error) {
	members, err := a.client.GetContactList(ctx)
	if err != nil {
		return nil, err
	}
	return OK(members), nil
}

// SyncMessage 同步消息
func (a *APIServer) SyncMessage(ctx context.Context, _ ginx.Empty) (*Result[[]*Message], error) {
	log.Ctx(ctx).Info().Msg("receive sync message request")
	messages := make([]*Message, 0)

	message, err := a.msgBuffer.Get(ctx, time.Second*25)
	if errors.Is(err, msgbuffer.ErrNoMessage) {
		return OK(messages), nil
	}
	if err != nil {
		return nil, err
	}
	messages = append(messages, message)
	return OK(messages), nil
}

func (a *APIServer) GetChatRoomDetail(ctx context.Context, req GetChatRoomInfoRequest) (*Result[*ChatRoomInfo], error) {
	info, err := a.client.GetChatRoomDetail(ctx, req.ChatRoomID)
	if err != nil {
		return nil, err
	}
	return OK(info), nil
}

func (a *APIServer) GetMemberFromChatRoom(ctx context.Context, req GetMemberFromChatRoomRequest) (*Result[[]*Profile], error) {
	members, err := a.client.GetMemberFromChatRoom(ctx, req.ChatRoomID)
	if err != nil {
		return nil, err
	}
	memberIds := strings.Split(members.Members, "^G")

	ctx, cancel := context.WithCancelCause(ctx)

	var eg errgroup.Group

	eg.SetLimit(runtime.NumCPU())

	result := make([]*Profile, len(memberIds))

	handler := func(index int, id string) func() error {
		return func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			profile, err := a.client.GetContactProfile(ctx, id)
			if err != nil {
				cancel(err)
				return err
			}
			result[index] = profile
			return nil
		}
	}
	for i, memberId := range memberIds {
		eg.Go(handler(i, memberId))
	}
	if err = eg.Wait(); err != nil {
		return nil, err
	}
	return OK(result), nil
}

func (a *APIServer) SendAtText(ctx context.Context, req SendAtTextRequest) (*Result[any], error) {
	if err := a.client.SendAtText(ctx, wxclient.SendAtTextOption{
		ChatRoomID: req.GroupID,
		WxIds:      req.AtList,
		Content:    req.Content,
	}); err != nil {
		return nil, err
	}
	return OK[any](nil), nil
}

func (a *APIServer) AddMemberToChatRoom(ctx context.Context, req AddMemberToChatRoomRequest) (*Result[any], error) {
	err := a.client.AddMemberIntoChatRoom(ctx, req.ChatRoomID, req.MemberIds)
	if err != nil {
		return nil, err
	}
	return OK[any](nil), nil
}

func (a *APIServer) InviteMemberToChatRoom(ctx context.Context, req InviteMemberToChatRoomRequest) (*Result[any], error) {
	err := a.client.InviteMemberToChatRoom(ctx, req.ChatRoomID, req.MemberIds)
	if err != nil {
		return nil, err
	}
	return OK[any](nil), nil
}

func (a *APIServer) ForwardMsg(ctx context.Context, req ForwardMsgRequest) (*Result[any], error) {
	err := a.client.ForwardMsg(ctx, req.WxID, req.MsgID)
	if err != nil {
		return nil, err
	}
	return OK[any](nil), nil
}

func (a *APIServer) startListen() error {
	// fixme: wine中无法根据service name进行dns解析，所以需要获取apiserver的ip
	addr, err := netutil.GetHostIP()
	if err != nil {
		return err
	}
	port := env.Name("MSG_LISTENER_PORT").IntOrElse(9999)

	{
		msgListener := &TCPMessageListener{Addr: ":" + strconv.Itoa(port)}
		// 定义消息处理行为，将获取到的消息塞进队列中
		var handler MessageHandlerFunc = func(message *Message) {
			_ = a.msgBuffer.Put(context.TODO(), message)
		}
		// 避免阻塞
		go func() {
			var stopReason error
			// 当 msgListener 停止之后 APIServer 也随之停止
			defer a.stop(stopReason)
			stopReason = msgListener.ListenAndServe(handler)
			log.Error().Err(stopReason).Msg("listen and serve message failed")
		}()
	}

	// 尝试去注册消息回调
	return a.client.HookSyncMsg(context.Background(), addr, port)
}

func (a *APIServer) Run(addr string) error {
	if err := a.startListen(); err != nil {
		return err
	}
	srv := &http.Server{
		Addr:    addr,
		Handler: registerAPIServer(a),
	}
	go a.checker.Check(a.ctx)
	return srv.ListenAndServe()
}

func New(client *wxclient.Client, msgBuffer msgbuffer.MessageBuffer) *APIServer {
	ctx, cancel := context.WithCancelCause(context.Background())
	srv := &APIServer{
		client:    client,
		msgBuffer: msgBuffer,
		ctx:       ctx,
		stop:      cancel,
	}
	srv.checker = &loginChecker{srv: srv, loopInterval: time.Second / 5}
	return srv
}

func Default() *APIServer {
	return New(wxclient.Default(), msgbuffer.Default())
}
