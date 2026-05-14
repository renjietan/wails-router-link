package fxLifeCycle

import (
	"context"
)

type FxLifecycle struct {
}

// OnStart 应用程序启动时执行
func (l *FxLifecycle) OnStart(context.Context) error {
	return nil
}

// OnStop 应用程序停止时执行
func (l *FxLifecycle) OnStop(context.Context) error {
	return nil
}

func NewFxLifeCycle() *FxLifecycle {
	return &FxLifecycle{}
}
