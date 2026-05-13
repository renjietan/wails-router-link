package middlewave

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"wails-router-link/service/core/component/logger"
	"wails-router-link/service/types"
)

// Http Request 全局错误处理
func ErrorHandlerMiddleWave(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.L().Errorf("errorHandler 异常捕获: %v", r)
			debug.PrintStack()
			c.JSON(http.StatusBadRequest, types.Response{Code: types.Failed, Message: types.ErrorMsg})
			c.Abort()
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
}
