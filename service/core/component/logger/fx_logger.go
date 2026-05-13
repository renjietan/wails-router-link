package logger

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx/fxevent"
)

type FxLogger struct {
	Logger        *logrus.Logger
	info          logrus.Fields // 需要打印的信息
	lock          sync.Mutex
	lifeCycleMame string //声明周期名称
}

func (l *FxLogger) LogEvent(event fxevent.Event) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.info = logrus.Fields{}
	switch e := event.(type) {
	case *fxevent.Provided:
		l.info["OutputTypeNames"] = e.OutputTypeNames
		l.info["ConstructorName"] = e.ConstructorName
		l.info["ModuleName"] = e.ModuleName
		l.info["ModuleTrace"] = e.ModuleTrace
		l.info["Private"] = e.Private
		l.info["StackTrace"] = e.StackTrace
		if e.Err != nil {
			l.info["Error"] = e.Err.Error()
			l.Error(event)
			return
		}
	case *fxevent.Invoking:
		l.info["ModuleName"] = e.ModuleName
		l.info["FunctionName"] = e.FunctionName
	case *fxevent.BeforeRun:
		l.info["ModuleName"] = e.ModuleName
		l.info["Name"] = e.Name
		l.info["Kind"] = e.Kind
	case *fxevent.Run:
		l.info["ModuleName"] = e.ModuleName
		l.info["Name"] = e.Name
		l.info["Kind"] = e.Kind
		l.info["Runtime"] = e.Runtime
		if e.Err != nil {
			l.info["Error"] = e.Err.Error()
			l.Error(event)
			return
		}
	case *fxevent.Invoked:
		l.info["FunctionName"] = e.FunctionName
		l.info["ModuleName"] = e.ModuleName
		l.info["Trace"] = e.Trace
		if e.Err != nil {
			l.info["Error"] = e.Err.Error()
			l.Error(event)
			return
		}
	case *fxevent.Started: // 应用启动完成
		if e.Err != nil {
			l.info["Error"] = e.Err.Error()
			l.Error(event)
			return
		}
	case *fxevent.OnStartExecuting: // 执行 OnStart 钩子前
		l.info["CallerName"] = e.CallerName
		l.info["FunctionName"] = e.FunctionName
	case *fxevent.OnStartExecuted: // 执行 OnStart 钩子后
		l.info["Runtime"] = e.Runtime
		l.info["CallerName"] = e.CallerName
		l.info["FunctionName"] = e.FunctionName
		l.info["Method"] = e.Method
		if e.Err != nil {
			l.info["Error"] = e.Err.Error()
			l.Error(event)
			return
		}
	case *fxevent.Stopped:
		l.Warn(event)
		return
	}
	l.Info(event)
}

func (l *FxLogger) info2string(event fxevent.Event) (name string) {
	val := reflect.TypeOf(event)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	name = val.Name()
	return
}

func (l *FxLogger) Info(event fxevent.Event) {
	name := l.info2string(event)
	name = fmt.Sprintf("✅ %s", name)
	l.Logger.WithFields(l.info).Info(name)
}

func (l *FxLogger) Debug(event fxevent.Event) {
	name := l.info2string(event)
	name = fmt.Sprintf("🔵 %s", name)
	l.Logger.WithFields(l.info).Debug(name)
}

func (l *FxLogger) Warn(event fxevent.Event) {
	name := l.info2string(event)
	name = fmt.Sprintf("️⚠️ %s", name)
	l.Logger.WithFields(l.info).Warn(name)
}

func (l *FxLogger) Error(event fxevent.Event) {
	name := l.info2string(event)
	name = fmt.Sprintf("❌ %s", name)
	l.Logger.WithFields(l.info).Error(name)
}

func (l *FxLogger) Fatal(event fxevent.Event) {
	name := l.info2string(event)
	name = fmt.Sprintf("❌ %s", name)
	l.Logger.WithFields(l.info).Fatal(name)
}

func (l *FxLogger) Panic(event fxevent.Event) {
	name := l.info2string(event)
	name = fmt.Sprintf("❌ %s", name)
	l.Logger.WithFields(l.info).Panic(name)
}
