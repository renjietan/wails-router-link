package middlewave

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/sirupsen/logrus"
	"wails-router-link/service/core/component/logger"
)

// GinLoggerMiddleWave Gin中间件：记录每个HTTP请求的日志
func GinLoggerMiddleWave() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		end := time.Now()
		latency := end.Sub(start)
		status := c.Writer.Status()
		method := c.Request.Method
		contentType := c.ContentType()
		fields := logrus.Fields{
			"package":     "GinLoggerMiddleWave",
			"method":      method,
			"contentType": contentType,
			"path":        c.Request.URL.Path,
			//"query":      c.Request.URL.RawQuery,
			"status":     status,
			"client_ip":  c.ClientIP(),
			"latency_ms": latency.Milliseconds(),
			"user_agent": c.Request.UserAgent(),
		}
		if len(c.Request.URL.Query()) > 0 {
			fields["Query"] = c.Request.URL.Query()
		}

		if strings.Contains(contentType, "multipart/form-data") {
			// 解析 multipart 表单，内存限制 32MB
			if err := c.Request.ParseMultipartForm(32 << 20); err == nil {
				// 获取普通表单字段
				if c.Request.MultipartForm != nil && len(c.Request.MultipartForm.Value) > 0 {
					fields["MultipartForm"] = c.Request.MultipartForm.Value
				}
				// 获取文件信息
				if c.Request.MultipartForm != nil && c.Request.MultipartForm.File != nil {
					for fieldName, fileHeaders := range c.Request.MultipartForm.File {
						for _, fh := range fileHeaders {
							fields["Files"] = map[string]interface{}{
								"FieldName": fieldName,
								"FileName":  fh.Filename,
								"Size":      fh.Size,
							}
						}
					}
				}
			}
		} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
			// 普通表单
			if err := c.Request.ParseForm(); err == nil {
				fields["PostForm"] = c.Request.PostForm
			}
		} else if strings.Contains(contentType, "application/json") {
			// JSON Body - 需要特殊处理，读取后重新放回
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil && len(bodyBytes) > 0 {
				// 将 body 重新写回，供后续处理使用
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// 尝试解析 JSON
				var jsonData interface{}
				if err := json.Unmarshal(bodyBytes, &jsonData); err == nil {
					fields["Body"] = jsonData
				} else {
					// 如果不是合法 JSON，记录原始字符串
					fields["Body"] = string(bodyBytes)
				}
			}
		}
		c.Next()
		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
		}
		// 在请求处理完成后补充路径参数
		params := make(map[string]string)
		for _, param := range c.Params {
			params[param.Key] = param.Value
		}
		if len(params) > 0 {
			fields["Params"] = params
		}
		logger.L().WithFields(fields).Info("HTTP Requeset: " + c.Request.Method + " " + c.Request.RequestURI)
		//WithFields(fields).Info("HTTP request")
	}
}
