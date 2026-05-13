package controller

import (
	"wails-router-link/service/api/service"
	"wails-router-link/service/core/component/appServer"
	"wails-router-link/service/utility/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConfigDetailsController struct {
	BaseController
	service *service.ConfigDetailsService
}

func NewConfigDetailController(
	app *appServer.AppServer,
	db *gorm.DB,
	service *service.ConfigDetailsService,
) *ConfigDetailsController {
	return &ConfigDetailsController{
		BaseController: BaseController{
			App: app,
			DB:  db,
		},
		service: service,
	}
}

func (configDetail *ConfigDetailsController) RegisterConfigDetailRouters() {
	group := configDetail.App.Engine.Group("/config_detail")
	group.POST(":id", configDetail.getDetailById)
}

// @Summary		 根据cId 获取 配置详情
// @Description  根据cId 获取 配置详情
// @Tags         ConfigDetail
// @Accept       json
// @Produce      json
// @Param        id	   path  	 int   true      "配置ID"
// @Success      200   {object}  types.Response  "上传成功"
// @Failure      400   {object}  types.Response  "请求参数错误"
// @Failure      500   {object}  types.Response  "服务器内部错误"
// @Router       /config_detail/{id} [post]
func (configDetial *ConfigDetailsController) getDetailById(c *gin.Context) {
	cId := c.Param("id")
	details := configDetial.service.GetById(cId)
	response.SUCCESS(c, details)
}
