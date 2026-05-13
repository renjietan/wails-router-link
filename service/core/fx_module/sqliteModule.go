package fx_module

import (
	"wails-router-link/service/api/service"
	"wails-router-link/service/core/component/sql_driver"

	"go.uber.org/fx"
)

var FXSQLiteModule = fx.Module("fx-sqlite-module",
	fx.Provide(sql_driver.NewSqliteDriver),
	// 自动同步 数据库 表
	fx.Provide(service.NewMigrationService),
	fx.Invoke(func(migrationService *service.MigrationService) {
		migrationService.StartMigrate()
	}),
)
