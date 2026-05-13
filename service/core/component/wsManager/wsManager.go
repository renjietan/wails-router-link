package wsManager

import (
	"fmt"
	"strings"

	"wails-router-link/service/core"
	"wails-router-link/service/types"
	"wails-router-link/service/utility"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

type WsManager struct {
	Server  *melody.Melody
	Clients map[string]*melody.Session
	config  *types.AppConfig
	Http    *core.AppServer
}

func NewWsManager(config *types.AppConfig, server *core.AppServer) *WsManager {
	config.WebSocket.Enable = true
	return &WsManager{
		Server:  melody.New(), // 创建 WS 示例
		Clients: make(map[string]*melody.Session),
		config:  config,
		Http:    server,
	}
}

func (ws *WsManager) InitRouter() {
	ws.Http.Engine.Any(ws.config.WebSocket.BaseUrl, func(context *gin.Context) {
		err := ws.Server.HandleRequest(context.Writer, context.Request)
		if err != nil {
			return
		}
	})
}

func (ws *WsManager) Listens() {
	ws.Server.HandleConnect(func(session *melody.Session) {
		paths := strings.Split(session.Request.RequestURI, "/")
		id := utility.Tern(len(paths) > 2, paths[2], "")
		key := ws.GenSessionKey()
		value := ws.GenClientKey(id)
		session.Set(key, value)
		ws.Clients[value] = session
		err := session.Write([]byte("connected"))
		if err != nil {
			return
		}
	})
	ws.Server.HandleDisconnect(func(s *melody.Session) {
		key := ws.GenSessionKey()
		delete(ws.Clients, key)
	})
	ws.Server.HandleMessage(func(session *melody.Session, bytes []byte) {
		msg := string(bytes)
		key := ws.GenSessionKey()
		id, _ := session.Get(key)
		fmt.Println("msg", msg)
		fmt.Println("id", id)
	})
}

func (ws *WsManager) GenSessionKey() string {
	return "id"
}

func (ws *WsManager) GenClientKey(id interface{}) string {
	return fmt.Sprintf("user-%s", id)
}
