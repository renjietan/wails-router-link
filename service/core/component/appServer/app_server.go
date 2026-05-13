package appServer

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"wails-router-link/service/core/middlewave"
	"wails-router-link/service/types"

	"github.com/gin-gonic/gin"
)

type AppServer struct {
	Config *types.AppConfig
	Engine *gin.Engine
	Server *http.Server
}

func NewAppServer(appConfig *types.AppConfig) *AppServer {
	//gin.DefaultWriter = logrus.StandardLogger().Out
	if !appConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
		gin.ForceConsoleColor()
	}
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		fmt.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	engine := gin.Default()
	addr := fmt.Sprintf("%s:%d", appConfig.Host, appConfig.Port)
	server := &http.Server{
		Addr:           addr,
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	res := &AppServer{
		Config: appConfig,
		Engine: engine,
		Server: server,
	}
	res.Middlewares()
	return res
}

func (s *AppServer) Middlewares() {
	//允许跨域请求 API
	s.Engine.Use(middlewave.CorsMiddleWave())
	// 静态资源
	s.Engine.Use(middlewave.StaticResourceMiddleWave())
	// 权限
	//s.Engine.Use(authorizeMiddleware(s, client))
	// 参数处理
	s.Engine.Use(middlewave.ParameterHandlerMiddleWave())
	// 错误处理
	s.Engine.Use(middlewave.ErrorHandlerMiddleWave)
	// 日志
	s.Engine.Use(middlewave.GinLoggerMiddleWave())
	//// 异常捕获: 捕获请求处理过程中发生的 panic
	s.Engine.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		res := gin.H{"error": "服务器内部错误", "detail": fmt.Sprintf("%v", err)}
		c.AbortWithStatusJSON(http.StatusInternalServerError, res)
	}))
	// 添加静态资源访问
	s.Engine.Static("/static", s.Config.StaticDir)
}

func (s *AppServer) Run() {
	// 注意：这里使用协程,
	// ListenAndServe 会阻塞 main 主协程，内部会启动一个循环，持续接受和处理连接
	go func() {
		if err := s.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
}
