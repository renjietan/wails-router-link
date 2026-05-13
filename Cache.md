package main

import (
"embed"
"log"
"wails-router-link/service"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
// Create an instance of the app structure
app := NewApp()
service.NewServiceApp()
// Create application with options
err := wails.Run(&options.App{
Title:             "wails-router-link",
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
OnStartup:        app.startup,
OnDomReady:       app.domReady,
OnBeforeClose:    app.beforeClose,
OnShutdown:       app.shutdown,
WindowStartState: options.Normal,
Bind: []interface{}{
app,
},
// Windows platform specific options
Windows: &windows.Options{
WebviewIsTransparent: false,
WindowIsTranslucent:  false,
DisableWindowIcon:    false,
// DisableFramelessWindowDecorations: false,
WebviewUserDataPath: "",
ZoomFactor:          1.0,
},
// Mac platform specific options
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
Title:   "wails-router-link",
Message: "",
Icon:    icon,
},
},
})

	if err != nil {
		log.Fatal(err)
	}
}

























package main

import (
"context"
"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {

	// Perform your setup here
	a.ctx = ctx
	runtime.EventsOn(ctx, "test-event", func(optionalData ...interface{}) {
		fmt.Println("数据接收：", optionalData[0].(string))
	})
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
return fmt.Sprintf("Hello %s, It's show time!", name)
}
