package web

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
	"xnps/pkg/database"
	"xnps/web/api"
	"xnps/web/service"
)

type Server struct {
	engin     *gin.Engine
	host      string
	kit       *service.Base
	CloseFlag chan struct{}
}

func NewServer() *Server {
	return &Server{}
}
func (w *Server) Close() {
	w.CloseFlag <- struct{}{}
}

func (w *Server) Start(host string, db *database.Driver) {
	w.engin = gin.Default()
	w.engin.Use(gin.Logger(), gin.Recovery())
	w.CloseFlag = make(chan struct{})
	w.kit = &service.Base{}
	w.kit = w.kit.Service(db)
	middle := NewMiddle(w.kit)
	xnps := w.engin.Group("/api/xnps")
	xnps.GET("/ping", api.Ping).
		POST("/login", middle.Login)
	//user
	userApi := api.NewUser(w.kit)
	userGroup := xnps.Group("/user", middle.AuthMiddle, middle.GetUser)
	userGroup.GET("/all", userApi.GetAllUser)
	//group
	groupApi := api.NewGroup(w.kit)
	group := xnps.Group("/group").Use(middle.AuthMiddle, middle.GetUser)
	group.GET("/all", groupApi.GetAll).
		POST("", groupApi.Create).
		POST("/filter", groupApi.GetByFilter).
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
	block := xnps.Group("/tunnel").Use(middle.AuthMiddle, middle.GetUser)
	block.GET("", blockApi.GetByIds).
		GET("/all", blockApi.GetAll).
		POST("/filter", blockApi.GetFilter).
		POST("/create", blockApi.Create).
		PUT("/:id", blockApi.Update).
		DELETE("/:id", blockApi.Delete)

	if err := w.engin.Run(host); err == nil {
		select {
		case <-w.CloseFlag:
			return
		}
	}
}

func (w *Server) ApiRouterRules(pre string) {

}

func (w *Server) RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(ctx *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			ctx.String(http.StatusForbidden, "rate limit...")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
