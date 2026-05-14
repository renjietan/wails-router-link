package fx_module

import (
	"context"
	"embed"
	"wails-router-link/service/core/component/fxLifeCycle"

	"go.uber.org/fx"
)

var FXLifeCycleModule = func(assets embed.FS, icon []byte) fx.Option {
	return fx.Module("fx-lifecycle-module",
		fx.Provide(fxLifeCycle.NewFxLifeCycle),
		fx.Invoke(func(lifecycle fx.Lifecycle, lc *fxLifeCycle.FxLifecycle) {
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
}
