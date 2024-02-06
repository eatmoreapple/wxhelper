package apiserver

import (
	"github.com/eatmoreapple/ginx"
)

func registerAPIServer(router *ginx.Router, server *APIServer) {
	router.GET(CheckLogin, ginx.G(server.CheckLogin).JSON())
	router.GET(GetUserInfo, ginx.G(server.GetUserInfo).JSON())
	router.GET(GetContactList, ginx.G(server.GetContactList).JSON())
	router.GET(SyncMessage, ginx.G(server.SyncMessage).JSON())
	router.POST(SendText, ginx.G(server.SendText).JSON())
	router.POST(SendImage, ginx.G(server.SendImage).JSON())
	router.POST(SendFile, ginx.G(server.SendFile).JSON())
	router.POST(GetChatRoomDetail, ginx.G(server.GetChatRoomDetail).JSON())
	router.POST(GetMemberFromChatRoom, ginx.G(server.GetMemberFromChatRoom).JSON())
}
