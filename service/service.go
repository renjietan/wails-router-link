package service

import (
	"context"
	"embed"
	"log"
	"os"
	"strconv"
	"wails-router-link/service/core"
	logger2 "wails-router-link/service/core/component/logger"
	"wails-router-link/service/core/fx_module"
	"wails-router-link/service/types"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

type AppLifecycle struct {
}

func NewServiceApp(assets embed.FS, icon []byte) {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config.toml"
	}
	debug, _ := strconv.ParseBool(os.Getenv("APP_DEBUG"))
	if !debug {
		defer func() {
			if err := recover(); err != nil {
				log.Fatal("❌ 生产环境 抛出异常(main.go): ", err)
			}
		}()
	}
	modules := getModules(assets, icon)
	fxApp := fx.New(
		fx.Provide(logger2.NewLogger),
		fx.WithLogger(func(l *logrus.Logger) fxevent.Logger {
			// TODO: 此处可能还需优化
			return &logger2.FxLogger{
				Logger: l,
			}
			//return fxevent.NopLogger
		}),
		fx.Provide(func() *types.AppConfig {
			config, err := core.LoadConfig(configFile)
			if err != nil {
				// 此处无法 使用 logrus
				// 因为 logger.NewLogger 中引入了 AppConfig， 此处 引入 logrus  会导致循环依赖引入
				log.Fatal("❌ 配置文件：读取失败")
			}
			config.ConfigPath = configFile
			config.Debug = debug
			if debug {
				_ = core.SaveConfig(config)
			}
			return config
		}),
		fx.Options(modules...),
	)
	go func() {
		if err := fxApp.Start(context.Background()); err != nil {
			log.Fatal("❌ fx启动失败", err)
		}
	}()
}

func getModules(assets embed.FS, icon []byte) []fx.Option {
	var fxOptions []fx.Option
	fxOptions = append(fxOptions, fx_module.FxGinModule)
	fxOptions = append(fxOptions, fx_module.FxGormConfigModule)
	fxOptions = append(fxOptions, fx_module.FXSwaggerModule)
	fxOptions = append(fxOptions, fx_module.FXSQLiteModule)
	fxOptions = append(fxOptions, fx_module.FXCronModule)
	fxOptions = append(fxOptions, fx_module.FxWsModule)
	fxOptions = append(fxOptions, fx_module.FXWailsAppModule(assets, icon))
	return fxOptions
}
