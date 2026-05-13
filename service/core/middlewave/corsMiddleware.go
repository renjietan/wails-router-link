package middlewave

import (
	"net/http"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
)

func CorsMiddleWave() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		// 设置允许的请求源
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}

		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		// 允许跨域设置可以返回其他子段，可以自定义字段
		c.Header("Access-Control-Allow-Headers", "Authorization, Body-Length, Body-Type, X-CSRF-Token,content-type")
		// 注意：允许浏览器（客户端）可以解析的头部
		c.Header("Access-Control-Expose-Headers", "Body-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		// 设置缓存时间
		c.Header("Access-Control-Max-Age", "172800")
		// 注意：允许客户端传递校验信息比如 cookie
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
		}

		defer func() {
			if err := recover(); err != nil {
				logger.Info("corsMiddleware中间件： 异常信息: %v", err)
			}
		}()

		c.Next()
	}
}
