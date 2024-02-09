package apiserver

import (
	"github.com/eatmoreapple/ginx"
	"github.com/gin-gonic/gin"
	"sync/atomic"
)

func registerAPIServer(router *ginx.Router, server *APIServer) {
	router.GET(CheckLogin, ginx.G(server.CheckLogin).JSON())

	{
		router.Use(loginRequired)
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

var loginRequired ginx.HandlerWrapper = func(ctx *gin.Context) error {
	apiserver := ctx.MustGet("apiserver").(*APIServer)
	if atomic.LoadInt32(&apiserver.status) != 1 {
		ctx.AbortWithStatusJSON(401, Result[any]{Code: resultCodeAuthErr, Msg: "not login"})
	}
	return nil
}
