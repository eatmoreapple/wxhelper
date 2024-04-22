package apiserver

import (
	"context"
	"github.com/eatmoreapple/ginx"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func initEngine(ctx context.Context) *gin.Engine {
	engine := gin.Default()
	engine.Use(func(c *gin.Context) { c.Request = c.Request.WithContext(ctx) })
	return engine
}

// activeRequired 要求用户登录后没有退出
func activeRequired(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		select {
		case <-ctx.Done():
			err := context.Cause(ctx)
			if err == nil {
				err = ctx.Err()
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, Result[any]{Code: resultCodeAuthErr, Msg: err.Error()})
		default:
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		}
	}
}

// loginRequired 要求用户登录
func loginRequired(authFunc func() bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !authFunc() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Result[any]{Code: resultCodeAuthErr, Msg: "not login"})
			return
		}
		c.Next()
	}
}

func registerAPIServer(server *APIServer) http.Handler {
	engine := initEngine(server.ctx)

	// global context middleware
	if onContext := server.OnContext; onContext != nil {
		engine.Use(func(c *gin.Context) {
			ctx := onContext(c.Request.Context())
			c.Request = c.Request.WithContext(ctx)
		})
	}

	pingRouter := ginx.NewRouter(engine)
	pingRouter.GET("/ping", ginx.G(server.Ping).String())

	router := ginx.NewRouter(engine)

	router.ErrorHandler = func(ctx *gin.Context, err error) {
		log.Ctx(ctx.Request.Context()).Error().Err(err).Msg("http error")
		ctx.JSON(http.StatusOK, Err[string](err.Error()))
	}

	engine.Use(activeRequired(server.ctx))

	checkLogin := ginx.G(server.CheckLogin).JSON()

	engine.GET(CheckLogin, func(c *gin.Context) {
		if err := checkLogin(c); err != nil {
			router.ErrorHandler(c, err)
		}
	})

	engine.Use(loginRequired(server.IsLogin))

	{
		router.GET(GetUserInfo, ginx.G(server.GetUserInfo).JSON())
		router.GET(GetContactList, ginx.G(server.GetContactList).JSON())
		router.GET(SyncMessage, ginx.G(server.SyncMessage).JSON())
		router.POST(SendText, ginx.G(server.SendText).JSON())
		router.POST(SendImage, ginx.G(server.SendImage).JSON())
		router.POST(SendFile, ginx.G(server.SendFile).JSON())
		router.POST(GetChatRoomDetail, ginx.G(server.GetChatRoomDetail).JSON())
		router.POST(GetMemberFromChatRoom, ginx.G(server.GetMemberFromChatRoom).JSON())
		router.POST(SendAtText, ginx.G(server.SendAtText).JSON())
		router.POST(AddMemberToChatRoom, ginx.G(server.AddMemberToChatRoom).JSON())
		router.POST(InviteMemberToChatRoom, ginx.G(server.InviteMemberToChatRoom).JSON())
		router.POST(ForwardMsg, ginx.G(server.ForwardMsg).JSON())
		router.POST(UploadFile, ginx.G(server.UploadFile).JSON())
	}
	return engine.Handler()
}
