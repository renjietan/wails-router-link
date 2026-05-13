package controller

import (
	"strings"
	"wails-router-link/service/core/component/appServer"

	"wails-router-link/service/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//var logger = logger2.GetLogger()

type BaseController struct {
	App *appServer.AppServer
	DB  *gorm.DB
}

func (h *BaseController) GetTrim(c *gin.Context, key string) string {
	return strings.TrimSpace(c.Query(key))
}

func (h *BaseController) PostInt(c *gin.Context, key string, defaultValue int) int {
	return utility.IntValue(c.PostForm(key), defaultValue)
}

func (h *BaseController) GetInt(c *gin.Context, key string, defaultValue int) int {
	return utility.IntValue(c.Query(key), defaultValue)
}

func (h *BaseController) GetFloat(c *gin.Context, key string) float64 {
	return utility.FloatValue(c.Query(key))
}
func (h *BaseController) PostFloat(c *gin.Context, key string) float64 {
	return utility.FloatValue(c.PostForm(key))
}

func (h *BaseController) GetBool(c *gin.Context, key string) bool {
	return utility.BoolValue(c.Query(key))
}
func (h *BaseController) PostBool(c *gin.Context, key string) bool {
	return utility.BoolValue(c.PostForm(key))
}
