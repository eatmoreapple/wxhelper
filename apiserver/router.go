package apiserver

import (
	"github.com/eatmoreapple/ginx"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync/atomic"
)

// todo 修改掉这里
func registerAPIServer(server *APIServer) {
	server.engine.Use(func(context *gin.Context) {
		context.Set("apiserver", server)
	})

	server.engine.GET(CheckLogin, func(context *gin.Context) {
		result, err := server.CheckLogin(context, struct{}{})
		if err != nil {
			context.JSON(http.StatusOK, Err[any](err.Error()))
		} else {
			context.JSON(http.StatusOK, result)
		}
	})
	{
		server.engine.Use(func(context *gin.Context) {
			apiserver := context.MustGet("apiserver").(*APIServer)
			if atomic.LoadInt32(&apiserver.status) != 1 {
				context.AbortWithStatusJSON(401, Result[any]{Code: resultCodeAuthErr, Msg: "not login"})
				return
			}
			context.Next()
		})
		router := ginx.NewRouter(server.engine)
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
}
