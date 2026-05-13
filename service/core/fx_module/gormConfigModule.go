package fx_module

import (
	"wails-router-link/service/core/component/sql_driver"

	"go.uber.org/fx"
)

var FxGormConfigModule = fx.Module("fx-gorm-config-module",
	fx.Provide(sql_driver.NewGormConfig),
)
