package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/web/gin/middleware"

	"github.com/gogoclouds/gogo/g"
	"github.com/gogoclouds/gogo/web/r"

	"github.com/gin-gonic/gin"
)

type RegisterHttpFn func(e *gin.Engine)

func RunHttpServer(exit <-chan struct{}, doneExit chan<- struct{}, addr string, register RegisterHttpFn) {
	e := gin.New()
	e.Use(gin.Logger()) // TODO -> zap.Logger
	e.Use(middleware.Recovery())
	e.Use(middleware.LoggerResponseFail())
	healthApi(e) // provide health API
	register(e)

	srv := &http.Server{Addr: addr, Handler: e}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-exit
	logger.Info("Shutting down http server...")
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer close(doneExit)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Http server forced to shutdown: ", err)
	}
	logger.Info("http server exiting")
}

// healthApi http check-up API
func healthApi(e *gin.Engine) {
	e.GET("/health", func(c *gin.Context) {
		appConf := g.Conf.App()
		msg := fmt.Sprintf("%s %s, is active", appConf.Name, appConf.Version)
		c.JSON(http.StatusOK, r.SuccessMsg(msg))
	})
}
