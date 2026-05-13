package fx_module

import (
	"context"
	"wails-router-link/service/core/component/appServer"
	"wails-router-link/service/core/component/fxLifeCycle"

	"go.uber.org/fx"
)

var FXLifeCycleModule = fx.Module("fx-lifecycle-module",
	fx.Provide(fxLifeCycle.NewFxLifeCycle),
	fx.Invoke(func(lifecycle fx.Lifecycle, lc *fxLifeCycle.FxLifecycle, server *appServer.AppServer) {
		lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return lc.OnStart(ctx)
			},
			OnStop: func(ctx context.Context) error {
				return lc.OnStop(ctx)
			},
		})
	}),
)
