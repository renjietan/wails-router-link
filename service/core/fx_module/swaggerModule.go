package fx_module

import (
	"go.uber.org/fx"
	"wails-router-link/service/core/component/swagger"
)

var FXSwaggerModule = fx.Module("fx-swagger-module",
	fx.Provide(swagger.NewSwaggerManager),
	fx.Invoke(func(sw *swagger.SwaggerManager) {
		sw.InitRouter()
	}),
)
