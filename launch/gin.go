package launch

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/whj1990/go-core/config"
	"github.com/whj1990/go-core/handler"
	"github.com/whj1990/go-core/middleware"
	"go.uber.org/zap"
)

type HttpRouter interface {
	SetRouter(*gin.Engine)
}

type RouterQuote struct {
	Routes []HttpRouter
}

func InitHttpServer(router ...HttpRouter) {
	app := gin.New()
	app.Use(middleware.ContextMiddleware())
	app.Use(gin.Recovery())
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app.GET("/ping", func(c *gin.Context) { handler.HandleSuccessResponse(c) })
	for _, r := range router {
		r.SetRouter(app)
	}

	srv := &http.Server{
		Handler: app,
		Addr:    net.JoinHostPort("0.0.0.0", config.GetNacosConfigData().HttpServer.Port),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				panic(err)
			} else {
				zap.L().Info("server graceful down")
			}
		}
	}()
	zap.L().Info("Start http server", zap.String("port", config.GetNacosConfigData().HttpServer.Port))
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx2, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx2); err != nil {
		panic(err)
	}
}
