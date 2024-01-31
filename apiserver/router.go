package apiserver

import (
	"github.com/eatmoreapple/ginx"
)

//func New() *gin.Engine {
//	engine := gin.Default()
//	router := ginx.NewRouter(engine)
//	server := &APIServer{}
//	registerAPIServer(router, server)
//	return engine
//}

func registerAPIServer(router *ginx.Router, server *APIServer) {
	router.POST(CheckLogin, ginx.G(server.CheckLogin).JSON())
	router.POST(GetUserInfo, ginx.G(server.GetUserInfo).JSON())
	router.POST(SendText, ginx.G(server.SendText).JSON())
	router.POST(GetContactList, ginx.G(server.GetContactList).JSON())
	router.POST(SyncMessage, ginx.G(server.SyncMessage).JSON())
}
