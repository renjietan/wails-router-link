package middlewave

import (
	"bytes"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"wails-router-link/service/utility"
)

// 中间件-参数处理
func ParameterHandlerMiddleWave() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET 参数处理
		params := c.Request.URL.Query()
		for key, values := range params {
			for i, value := range values {
				params[key][i] = strings.TrimSpace(value)
			}
		}
		// update get parameters
		c.Request.URL.RawQuery = params.Encode()
		// skip file upload requests
		contentType := c.Request.Header.Get("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {
			c.Next()
			return
		}

		if strings.Contains(contentType, "application/json") {
			// process POST JSON request body
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.Next()
				return
			}

			// 还原请求体
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			// 将请求体解析为 JSON
			var jsonData map[string]interface{}
			if err := c.ShouldBindJSON(&jsonData); err != nil {
				c.Next()
				return
			}

			// 对 JSON 数据中的字符串值去除两端空格
			trimJSONStrings(jsonData)
			// 更新请求体
			c.Request.Body = io.NopCloser(bytes.NewBufferString(utility.JsonEncode(jsonData)))
		}
		c.Next()
	}
}

func trimJSONStrings(data interface{}) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			switch valueType := value.(type) {
			case string:
				v[key] = strings.TrimSpace(valueType)
			case map[string]interface{}, []interface{}:
				trimJSONStrings(value)
			}
		}
	case []interface{}:
		for i, value := range v {
			switch valueType := value.(type) {
			case string:
				v[i] = strings.TrimSpace(valueType)
			case map[string]interface{}, []interface{}:
				trimJSONStrings(value)
			}
		}
	}
}
