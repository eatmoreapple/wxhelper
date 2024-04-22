package apiserver

import (
	"context"
	"errors"
	"github.com/eatmoreapple/env"
	"github.com/eatmoreapple/ginx"
	"github.com/eatmoreapple/wxhelper/apiserver/internal/filemerger"
	"github.com/eatmoreapple/wxhelper/apiserver/internal/msgbuffer"
	. "github.com/eatmoreapple/wxhelper/internal/models"
	"github.com/eatmoreapple/wxhelper/internal/wxclient"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var ErrLogout = errors.New("logout")

// APIServer 用来屏蔽微信的接口
type APIServer struct {
	client            *wxclient.Client
	msgBuffer         msgbuffer.MessageBuffer
	fileMergerFactory filemerger.Factory
	status            int32
	ctx               context.Context
	stop              context.CancelCauseFunc
	checker           Checker
	OnContext         func(context.Context) context.Context
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

type SendImageRequest struct {
	To    string `json:"to"`
	Image string `json:"image"`
}

func (a *APIServer) SendImage(ctx context.Context, req SendImageRequest) (*Result[any], error) {
	err := a.client.SendImage(ctx, req.To, req.Image)
	if err != nil {
		return nil, err
	}
	return OK[any](nil), nil
}

type SendFileRequest struct {
	To   string `json:"to"`
	File string `json:"file"`
}

func (a *APIServer) SendFile(ctx context.Context, req SendFileRequest) (*Result[any], error) {
	err := a.client.SendFile(ctx, req.To, req.File)
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

type GetChatRoomInfoRequest struct {
	ChatRoomID string `json:"chatRoomId"`
}

func (a *APIServer) GetChatRoomDetail(ctx context.Context, req GetChatRoomInfoRequest) (*Result[*ChatRoomInfo], error) {
	info, err := a.client.GetChatRoomDetail(ctx, req.ChatRoomID)
	if err != nil {
		return nil, err
	}
	return OK(info), nil
}

type GetMemberFromChatRoomRequest struct {
	ChatRoomID string `json:"chatRoomId"`
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

type SendAtTextRequest struct {
	GroupID string   `json:"groupId"`
	AtList  []string `json:"atList"`
	Content string   `json:"content"`
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

type AddMemberToChatRoomRequest struct {
	ChatRoomID string   `json:"chatRoomId"`
	MemberIds  []string `json:"memberIds"`
}

func (a *APIServer) AddMemberToChatRoom(ctx context.Context, req AddMemberToChatRoomRequest) (*Result[any], error) {
	err := a.client.AddMemberIntoChatRoom(ctx, req.ChatRoomID, req.MemberIds)
	if err != nil {
		return nil, err
	}
	return OK[any](nil), nil
}

type InviteMemberToChatRoomRequest struct {
	ChatRoomID string   `json:"chatRoomId"`
	MemberIds  []string `json:"memberIds"`
}

func (a *APIServer) InviteMemberToChatRoom(ctx context.Context, req InviteMemberToChatRoomRequest) (*Result[any], error) {
	err := a.client.InviteMemberToChatRoom(ctx, req.ChatRoomID, req.MemberIds)
	if err != nil {
		return nil, err
	}
	return OK[any](nil), nil
}

type ForwardMsgRequest struct {
	WxID  string `json:"wxid"`
	MsgID string `json:"msgId"`
}

func (a *APIServer) ForwardMsg(ctx context.Context, req ForwardMsgRequest) (*Result[any], error) {
	err := a.client.ForwardMsg(ctx, req.WxID, req.MsgID)
	if err != nil {
		return nil, err
	}
	return OK[any](nil), nil
}

type UploadRequest struct {
	Filename string        `form:"filename"`
	FileHash string        `form:"fileHash"`
	Chunks   int           `form:"chunks"`
	Chunk    int           `form:"chunk"`
	Content  io.ReadCloser `form:"-"`
}

func (a *UploadRequest) FromContext(ctx *gin.Context) error {
	if err := ctx.ShouldBind(a); err != nil {
		return err
	}
	reader, _, err := ctx.Request.FormFile("file")
	if err != nil {
		return err
	}
	a.Content = reader
	return nil
}

func (a *APIServer) UploadFile(ctx context.Context, req UploadRequest) (*Result[string], error) {
	// 保存上传的文件
	// 保存到本地
	file, err := os.CreateTemp("", "upload")
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	defer func() { _ = req.Content.Close() }()

	// save the file
	if _, err = io.Copy(file, req.Content); err != nil {
		return nil, err
	}
	key := req.Filename + ":" + req.FileHash

	fileMerger, err := a.fileMergerFactory.New(key)
	if err != nil {
		return nil, err
	}
	if err = fileMerger.Add(ctx, file.Name()); err != nil {
		return nil, err
	}
	var filename string

	// if it is the last chunk, merge the file
	if req.Chunk+1 == req.Chunks {
		// merge the file
		filename, err = fileMerger.Merge(ctx)
		if err != nil {
			return nil, err
		}
	}
	return OK[string](filename), nil
}

func (a *APIServer) startListen() error {
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
			log.Ctx(a.ctx).Error().Err(stopReason).Msg("listen and serve message failed")
		}()
	}
	// 尝试去注册消息回调
	// 已经在一个容器内了，直接用localhost
	return a.client.HookSyncMsg(a.ctx, "localhost", port)
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

func New(client *wxclient.Client, fileMergerFactory filemerger.Factory, msgBuffer msgbuffer.MessageBuffer) *APIServer {
	ctx, cancel := context.WithCancelCause(context.Background())
	srv := &APIServer{
		client:            client,
		msgBuffer:         msgBuffer,
		fileMergerFactory: fileMergerFactory,
		ctx:               ctx,
		stop:              cancel,
	}
	srv.checker = &loginChecker{srv: srv, loopInterval: time.Second / 5}
	return srv
}

func Default() *APIServer {
	return New(wxclient.Default(), filemerger.DefaultFactory(), msgbuffer.Default())
}
