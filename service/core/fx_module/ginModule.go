package fx_module

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"wails-router-link/service/core"
)

var FxGinModule = fx.Module("fx-gin-module",
	fx.Provide(core.NewAppServer),
	FXApiModule,
	fx.Invoke(func(appserver *core.AppServer, db *gorm.DB, log *logrus.Logger) {
		appserver.Run()
	}),
)
