package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"time"
	"tunpx/pkg/database"
	"tunpx/web/api"
	"tunpx/web/service"
)

type Server struct {
	engin     *gin.Engine
	ExCode    int
	host      string
	kit       *service.Base
	CloseFlag chan struct{}
}

func New() *Server {
	return &Server{
		engin:     gin.Default(),
		CloseFlag: make(chan struct{}),
		kit:       &service.Base{},
	}
}
func (w *Server) Close() {
	w.CloseFlag <- struct{}{}

}

func (w *Server) InitSys(host string, db *database.Driver) (err error) {
	// w.engin = gin.Default()
	w.engin.Use(gin.Logger(), gin.Recovery())
	w.kit = w.kit.Service(db)
	system := api.NewSystem(w.kit)
	w.engin.GET("/static/system/init", system.StaticInit).
		GET("/static/system/success", system.StaticSuccess).
		POST("/api/system", system.Init, func(ctx *gin.Context) {
			if !ctx.IsAborted() {
				go func() {
					time.Sleep(time.Millisecond * 500)
					w.Close()
				}()
				w.ExCode = 0
			}
		})
	w.host = host
	return w.listen(host)
}

func (w *Server) listen(host string) (err error) {
	sev := http.Server{Addr: host, Handler: w.engin}
	go func() {
		if err = sev.ListenAndServe(); err != nil {
			slog.Error("listen web server", "error", err)
		}
	}()

	select {
	case <-w.CloseFlag:
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err = sev.Shutdown(ctx); err != nil {
			slog.Error("init web server close failed", "error", err)
		} else {
			slog.Info("init web server close succeed")
		}
		return err
	}

}

func (w *Server) Start(host string, db *database.Driver) {
	w.engin.Use(gin.Logger(), gin.Recovery())
	w.kit = w.kit.Service(db)
	middle := NewMiddle(w.kit)

	xnps := w.engin.Group("/api/tunpxs")
	xnps.GET("/ping", middle.RateLimitMiddle(time.Second, 100, 10), api.Ping).
		POST("/login", middle.RateLimitMiddle(time.Second, 100, 10), middle.Login)
	// user
	userApi := api.NewUser(w.kit)
	userGroup := xnps.Group("/user", middle.AuthMiddle, middle.GetUser)
	userGroup.GET("/all", userApi.GetAllUser)
	// group
	groupApi := api.NewGroup(w.kit)
	group := xnps.Group("/group").Use(middle.AuthMiddle, middle.GetUser)
	group.GET("", groupApi.GetAll).
		POST("", groupApi.Create).
		PUT("/:id", groupApi.Update).
		DELETE("/:id", groupApi.Delete)

	clientApi := api.NewClient(w.kit)
	client := xnps.Group("/client").Use(middle.AuthMiddle, middle.GetUser)
	client.GET("", clientApi.GetByIds).
		GET("/all", clientApi.GetAll).
		POST("/filter", clientApi.GetFilter).
		POST("/create", clientApi.Create).
		PUT("/:id", clientApi.Update).
		DELETE("/:id", clientApi.Delete)

	tunnelApi := api.NewTunnel(w.kit)
	tunnel := xnps.Group("/tunnel").Use(middle.AuthMiddle, middle.GetUser)
	tunnel.GET("", tunnelApi.GetByIds).
		GET("/all", tunnelApi.GetAll).
		POST("/filter", tunnelApi.GetFilter).
		POST("/create", tunnelApi.Create).
		PUT("/:id", tunnelApi.Update).
		DELETE("/:id", tunnelApi.Delete)

	blockApi := api.NewBlock(w.kit)
	block := xnps.Group("/block").Use(middle.AuthMiddle, middle.GetUser)
	block.GET("", blockApi.GetByIds).
		GET("/all", blockApi.GetAll).
		POST("/filter", blockApi.GetFilter).
		POST("/create", blockApi.Create).
		PUT("/:id", blockApi.Update).
		DELETE("/:id", blockApi.Delete)
	go func() {
		err := w.engin.Run(host)
		if err != nil {
			slog.Error("start web error", "error", err)
		}
	}()
	select {
	case <-w.CloseFlag:
		return
	}
}
