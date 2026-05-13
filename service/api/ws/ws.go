package ws

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"wails-router-link/service/enum"
	"wails-router-link/service/utility"
)

type ClientInfo struct {
	Name   string
	Client *websocket.Conn
	info   string
}

type WebSocketManager struct {
	// 所有连接的客户端
	Clients map[*websocket.Conn]ClientInfo
	// 用于广播消息的通道
	Broadcast_Msg_Chan chan []byte
	// 用于注册新客户端
	Register_Chan chan *websocket.Conn
	// 用于注销客户端
	Unregister_Chan chan *websocket.Conn
	mu              sync.RWMutex
	upgrader        websocket.Upgrader
	closed          atomic.Bool
	l               *logrus.Logger
}

// NewWebSocketManager 创建并初始化 WebSocket 管理器
func NewWebSocketManager(l *logrus.Logger) *WebSocketManager {
	return &WebSocketManager{
		Clients:            make(map[*websocket.Conn]ClientInfo),
		Broadcast_Msg_Chan: make(chan []byte, 256),
		Register_Chan:      make(chan *websocket.Conn),
		Unregister_Chan:    make(chan *websocket.Conn),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// 允许跨域
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		l: l,
	}
}

// HandleWebSocket 客户端-新增
func (m *WebSocketManager) HandleWebSocket(c *gin.Context) {
	conn, err := m.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		m.l.Error("WebSocket 升级失败: %v", err)
		return
	}
	// 注册新客户端
	m.Register_Chan <- conn

	// 启动读取协程
	go m.read(conn)
}

// 消息发送
func (m *WebSocketManager) SendToClient(conn *websocket.Conn, message []byte) error {
	m.mu.RLock()
	_, exists := m.Clients[conn]
	m.mu.RUnlock()
	if !exists {
		m.l.Error("客户端连接不存在")
		return fmt.Errorf("客户端连接不存在")
	}
	return conn.WriteMessage(websocket.TextMessage, message)
}

// 广播
func (m *WebSocketManager) Broadcast(message []byte) {
	m.Broadcast_Msg_Chan <- message
}

// GET 客户端数量
func (m *WebSocketManager) GetClientCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.Clients)
}

// 断开连接
func (m *WebSocketManager) Close() {
	if !m.closed.CompareAndSwap(false, true) {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	// 关闭所有客户端连接
	for conn := range m.Clients {
		err := conn.Close()
		if err != nil {
			m.l.Error(err.Error())
		}
		delete(m.Clients, conn)
	}
	close(m.Broadcast_Msg_Chan)
	close(m.Register_Chan)
	close(m.Unregister_Chan)
}

// 主循环
func (m *WebSocketManager) Run() {
	for {
		select {
		case conn := <-m.Register_Chan:
			// 注册新客户端
			m.mu.Lock()
			m.Clients[conn] = ClientInfo{
				Name:   conn.RemoteAddr().String(),
				Client: conn,
				info:   "test",
			}
			m.mu.Unlock()
			log.Printf("新客户端连接，当前连接数: %d", len(m.Clients))
		case conn := <-m.Unregister_Chan:
			// 注销客户端
			m.mu.Lock()
			if _, ok := m.Clients[conn]; ok {
				delete(m.Clients, conn)
				conn.Close()
				log.Printf("客户端断开连接，当前连接数: %d", len(m.Clients))
			}
			m.mu.Unlock()

		case message := <-m.Broadcast_Msg_Chan:
			// 广播消息给所有客户端
			m.mu.RLock()
			for conn := range m.Clients {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("发送消息失败: %v", err)
					// 如果发送失败，将连接加入注销队列
					select {
					case m.Unregister_Chan <- conn:
					default:
					}
				}
			}
			m.mu.RUnlock()
		}
	}
}

// 读取客户端消息的协程（没写完 - 数据处理）
func (m *WebSocketManager) read(conn *websocket.Conn) {
	defer func() {
		m.Unregister_Chan <- conn
		conn.Close()
	}()

	// 设置读取超时和 pong 处理
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	for {
		// 读取消息
		messageType, message, err := conn.ReadMessage()
		fmt.Println("消息类型：", messageType)
		// 检查是否是 正常规避 或 读取超时
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket 错误: %v", err)
			}
			break
		}
		// 任意消息都刷新读取超时，避免仅依赖控制帧 pong
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		// 下面 可以开始写我的逻辑了
		log.Printf("收到客户端消息: %s, \n", string(message))

		response := fmt.Sprintf("服务器收到: %s", string(message))
		// 向客户端发送消息
		if err := m.SendToClient(conn, []byte(response)); err != nil {
			log.Printf("发送回显消息失败: %v", err)
			break
		}
	}
}

// TODO：发送时 的数据处理（没写完）
func (m *WebSocketManager) handleRecv(message string) (s string, err any) {
	recJson := map[string]interface{}{}
	err = utility.Interface2Interface(string(message), &m)
	if err != nil {
		return "", err
	}
	e := recJson["Event"]
	dataStr := recJson["Data"]
	switch e {
	case enum.WS_EVENT_PING:

		str := utility.Interface2String(&map[string]any{
			"Event": enum.WS_EVENT_PONG,
			"Data":  dataStr,
			"type":  enum.WS_TYPE_SERVER,
		})
		return str, nil
	default:
		return utility.Interface2String(&map[string]any{
			"Event": enum.WS_EVENT_UNKNOWN,
			"Data":  dataStr,
			"type":  enum.WS_TYPE_SERVER,
		}), nil
	}
}
