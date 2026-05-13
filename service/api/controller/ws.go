package controller

import (
	"github.com/gin-gonic/gin"
	"wails-router-link/service/api/ws"
)

// WsClientCountResponse 用于 swagger 展示当前连接数响应体.
type WsClientCountResponse struct {
	ClientCount int `json:"client_count" example:"3"`
}

// BroadcastRequest 用于 swagger 展示广播请求体.
type BroadcastRequest struct {
	Msg string `json:"msg" example:"hello"`
}

// StatusResponse 通用状态响应.
type StatusResponse struct {
	Status string `json:"status" example:"ok"`
}

// ErrorResponse 通用错误响应.
type ErrorResponse struct {
	Error  string `json:"error" example:"需要字段 msg"`
	Detail string `json:"detail,omitempty" example:"dial udp: ..."`
}

// GetWsClientCount godoc
//
//	@Summary	获取当前 WebSocket 连接数
//	@Tags		websocket
//	@Produce	json
//	@Success	200	{object}	WsClientCountResponse
//	@Router		/ws/count [get]
func GetWsClientCount(c *gin.Context) {
	ws := c.MustGet("ws").(*ws.WebSocketManager)
	count := ws.GetClientCount()
	c.JSON(200, gin.H{
		"client_count": count,
	})
}

// BroadcastWsMessage godoc
//
//	@Summary	广播 WebSocket 消息
//	@Tags		websocket
//	@Accept		json
//	@Produce	json
//	@Param		data	body		BroadcastRequest	true	"消息体"
//	@Success	200		{object}	StatusResponse
//	@Failure	400		{object}	ErrorResponse
//	@Router		/ws/broadcast [post]
func BroadcastWsMessage(c *gin.Context) {
	ws := c.MustGet("ws").(*ws.WebSocketManager)
	var req BroadcastRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Msg == "" {
		c.JSON(400, gin.H{"error": "需要字段 msg"})
		return
	}

	ws.Broadcast([]byte(req.Msg))
	c.JSON(200, gin.H{"status": "ok"})
}
