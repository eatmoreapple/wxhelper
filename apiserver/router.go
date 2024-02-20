package apiserver

import (
	stdcontext "context"
	"github.com/eatmoreapple/ginx"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"net/http"
)

func initEngine() *gin.Engine {
	engine := gin.Default()
	engine.Use(func(c *gin.Context) { c.Request = c.Request.WithContext(log.Logger.WithContext(c.Request.Context())) })
	return engine
}

// activeRequired 要求用户登录后没有退出
func activeRequired(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		select {
		case <-ctx.Done():
			err := stdcontext.Cause(ctx)
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
	engine := initEngine()

	engine.Use(activeRequired(server.ctx))

	engine.Use(loginRequired(server.IsLogin))

	{
		router := ginx.NewRouter(engine)
		router.GET(CheckLogin, ginx.G(server.CheckLogin).JSON())
		router.GET(GetUserInfo, ginx.G(server.GetUserInfo).JSON())
		router.GET(GetContactList, ginx.G(server.GetContactList).JSON())
		router.GET(SyncMessage, ginx.G(server.SyncMessage).JSON())
		router.POST(SendText, ginx.G(server.SendText).JSON())
		router.POST(SendImage, ginx.G(server.SendImage).JSON())
		router.POST(SendFile, ginx.G(server.SendFile).JSON())
		router.POST(GetChatRoomDetail, ginx.G(server.GetChatRoomDetail).JSON())
		router.POST(GetMemberFromChatRoom, ginx.G(server.GetMemberFromChatRoom).JSON())
		router.POST(SendAtText, ginx.G(server.SendAtText).JSON())
	}
	return engine.Handler()
}
