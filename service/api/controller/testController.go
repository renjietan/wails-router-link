package controller

import (
	"fmt"
	"net/http"
	"wails-router-link/service/core/component/appServer"

	dto "wails-router-link/service/api/DTO"

	"github.com/gin-gonic/gin"
)

type TestController struct {
	BaseController
}

func NewTestController(
	app *appServer.AppServer,
) *ConfigController {
	return &ConfigController{
		BaseController: BaseController{
			App: app,
		},
	}
}

func (t *TestController) RegisterTestRouters() {
	group := t.App.Engine.Group("/test")
	group.GET("list", t.Form)
	//group.POST("insert", config.insert)
	//group.PATCH("insertMany", config.insertMany)
	//group.POST("upload", config.upload)
	//group.PUT("update/:id", config.update)
	//group.PUT("updates", config.updates)
	//group.DELETE("delete/:id", config.delete)
	//group.GET("test", config.test)
}

// @Summary form
// @Description form
// @Tags Test
// @Accept json
// @Produce json
// @Param request body dto.TestFormDTO true "参数"
// @Router /test/form [post]
func (config *TestController) Form(c *gin.Context) {
	var params dto.TestFormDTO
	if err := c.ShouldBind(&params); err != nil {
		fmt.Println(err)
	}
	c.AsciiJSON(http.StatusOK, params)
}
