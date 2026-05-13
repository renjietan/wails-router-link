package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"wails-router-link/service/utility"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type LinePerFieldFormatter struct {
	prefixed.TextFormatter
}

// Format 实现 logrus.Formatter 接口
func (f *LinePerFieldFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buf bytes.Buffer
	timestamp := entry.Time.Format(f.TimestampFormat)
	level := strings.ToUpper(entry.Level.String())

	prefix := ""
	if p, ok := entry.Data["prefix"]; ok {
		prefix = fmt.Sprintf("%v", p)
	}

	header := fmt.Sprintf("[%s] %-5s %s %s\n",
		timestamp,
		level,
		prefix,
		entry.Message,
	)
	buf.WriteString(header)

	for key, value := range entry.Data {
		if key == "prefix" {
			continue
		}
		fieldLine := fmt.Sprintf("  %s: %v\n", key, value)
		buf.WriteString(fieldLine)
	}

	if entry.HasCaller() {
		callerLine := fmt.Sprintf("  caller: %s:%d %s\n",
			entry.Caller.File, entry.Caller.Line, entry.Caller.Function)
		buf.WriteString(callerLine)
	}

	return buf.Bytes(), nil
}

func BeautifyJsonStr(obj map[string]interface{}) string {
	rawJSON := utility.Interface2String(&obj)
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(rawJSON), "", "  ") // 前缀为空，缩进为两个空格
	if err != nil {
		fmt.Println("BeautifyJsonStr error:", err)
	}
	return prettyJSON.String()
}
