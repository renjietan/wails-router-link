package fx_module

import (
	"wails-router-link/service/core/component/appServer"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var FxGinModule = fx.Module("fx-gin-module",
	fx.Provide(appServer.NewAppServer),
	FXApiModule,
	fx.Invoke(func(appserver *appServer.AppServer, db *gorm.DB, log *logrus.Logger) {
		appserver.Run()
	}),
)
