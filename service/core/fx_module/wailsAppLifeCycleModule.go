package fx_module

import (
	"embed"
	"log"
	"wails-router-link/service/core/component/wailsAppLifeCycle"
	"wails-router-link/service/types"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"go.uber.org/fx"
)

// FXWailsAppLifeCycleModule 暂时废弃
var FXWailsAppLifeCycleModule = func(assets embed.FS, icon []byte) fx.Option {
	return fx.Module("fx-wailsApp-module",
		fx.Provide(wailsAppLifeCycle.NewApp),
		fx.Invoke(func(wailsApp *wailsAppLifeCycle.App, appConfig *types.AppConfig) {

		}),
	)
}

var RunWails = func(wailsApp *wailsAppLifeCycle.App, appConfig *types.AppConfig, assets embed.FS, icon []byte) {
	err := wails.Run(&options.App{
		Title:             appConfig.AppName,
		Width:             1024,
		Height:            768,
		MinWidth:          1024,
		MinHeight:         768,
		MaxWidth:          1280,
		MaxHeight:         800,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Menu:             nil,
		Logger:           nil,
		LogLevel:         logger.DEBUG,
		OnStartup:        wailsApp.Startup,
		OnDomReady:       wailsApp.DomReady,
		OnBeforeClose:    wailsApp.BeforeClose,
		OnShutdown:       wailsApp.Shutdown,
		WindowStartState: options.Normal,
		Bind: []interface{}{
			wailsApp,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
			ZoomFactor:          1.0,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   appConfig.AppName,
				Message: "",
				Icon:    icon,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
