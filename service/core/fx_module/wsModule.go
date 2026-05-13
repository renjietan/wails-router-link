package fx_module

import (
	"wails-router-link/service/core/component/wsManager"

	"go.uber.org/fx"
)

var FxWsModule = fx.Module("fx-ws-module",
	//fx.Provide(ws.NewWebSocketManager),
	//fx.Invoke(func(ws *ws.WebSocketManager, appserver *core.AppServer) {
	//	//appserver.Engine.Use(func(c *gin.Context) {
	//	//	c.Set("ws", ws.HandleWebSocket)
	//	//	c.Next()
	//	//})
	//}),
	fx.Provide(wsManager.NewWsManager),
	fx.Invoke(func(ws *wsManager.WsManager) {
		ws.InitRouter()
		ws.Listens()
	}),
)
