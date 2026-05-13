package swagger

import (
	"fmt"
	"wails-router-link/service/core/component/appServer"

	"wails-router-link/service/docs"
	"wails-router-link/service/types"
	"wails-router-link/service/utility/response"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SwaggerManager struct {
	Server      *appServer.AppServer
	HandlerFunc gin.HandlerFunc
	Config      *types.AppConfig
}

func NewSwaggerManager(appserver *appServer.AppServer, appConfig *types.AppConfig) *SwaggerManager {
	var swaggerHandler gin.HandlerFunc
	if appConfig.Debug == false {
		swaggerHandler = func(c *gin.Context) {
			response.ERROR(c, "swagger 中间件已被禁用")
		}
	} else {
		url := fmt.Sprintf("http://%s:%d/swagger/doc.json", appConfig.Host, appConfig.Port)
		swaggerHandler = ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.URL(url),                     // 1. 显示指定 doc.json 的 URL
			ginSwagger.DefaultModelsExpandDepth(-1), // 2. 设置模型默认展开深度，-1 为完全隐藏
			ginSwagger.DocExpansion("none"),         // 3. 设置文档展开模式，'none' 为全部折叠
			ginSwagger.DeepLinking(true),            // 4. 开启深度链接
			ginSwagger.PersistAuthorization(true),   // 5. 持久化认证信息（如 API Key）
		)
	}
	appConfig.Swagger.Enable = true
	return &SwaggerManager{
		Server:      appserver,
		Config:      appConfig,
		HandlerFunc: swaggerHandler,
	}
}

func (sc *SwaggerManager) InitRouter() {
	docs.SwaggerInfo.Title = sc.Config.AppName
	docs.SwaggerInfo.Description = fmt.Sprintf("%s Swagger Document", sc.Config.AppName)
	docs.SwaggerInfo.Version = sc.Config.Version
	sc.Server.Engine.GET("/swagger/*any", sc.HandlerFunc)
}
