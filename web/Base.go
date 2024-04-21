package web

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
	"xnps/pkg/crypt"
	"xnps/pkg/database"
	"xnps/web/api"
	"xnps/web/dto"
	"xnps/web/service"
)

type Server struct {
	g         *gin.Engine
	CloseFlag chan struct{}
}

func NewServer() *Server {
	return &Server{}
}

func (w *Server) Start(host string, db *database.Driver) {
	w.g = gin.Default()
	w.g.Use(gin.Logger(), gin.Recovery())
	xnps := w.g.Group("/api/xnps")
	xnps.GET("/ping", api.Ping).
		GET("/doLogin", w.DoLogin).
		POST("/login", w.Login)
	kit := &service.Base{}
	kit = kit.Service(db)
	middle := NewMiddle(kit)

	userApi := service.NewAuthUser(kit)
	userGroup := xnps.Group("/user", middle.AuthMiddle, middle.GetUser)
	userGroup.GET("/all")

	r := xnps.Group("/group").Use(middle.AuthMiddle, middle.GetUser)

	//分组管理
	r.POST("/group/get/all", api.GetAllGroup)
	r.POST("/group/get/condition", api.GetGroupByCondition)
	r.POST("/group/add/condition", api.AddGroup)
	r.POST("/group/update/one", api.EditGroup)
	r.POST("/group/delete/one", api.DelGroup)
	//终端相关
	r.POST("/client/get/all", api.GetAllClient)
	r.POST("/client/get/condition", api.GetGroupByCondition)
	r.POST("/client/add/condition", api.AddClient)
	r.POST("/client/update/one", api.EditClient)
	r.POST("/client/delete/one", api.DelClient)
	//隧道相关
	r.POST("/tunnel/get/all", api.GetAllTunnel)
	r.POST("/tunnel/get/condition", api.GetTunnelByCondition)
	r.POST("/tunnel/add/condition", api.AddTunnel)
	r.POST("/tunnel/update/one", api.EditTunnel)
	r.POST("/tunnel/delete/one", api.DelTunnel)
	//黑名单相关
	r.POST("/block/get/all", api.GetAllBlockList)
	r.POST("/block/get/condition", api.GetBlockListByCondition)
	r.POST("/block/add/condition", api.AddFBlockList)
	r.POST("/block/update/one", api.EditBlockList)
	r.POST("/block/delete/one", api.DelBlockList)

	w.ApiRouterRules("/Api")
	if err := w.g.Run(host); err == nil {
		select {
		case <-w.CloseFlag:

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

// Login 登录信息
func (w *Server) Login(ctx *gin.Context) {
	var login dto.Login
	var token string
	var err error

	if err = ctx.ShouldBindJSON(&login); err == nil {
		if Salt, ok := w.salt.GetString(login.Username); ok {
			var passwd string
			if passwd, err = service.GetDb().GetPasswdByUser(login.Username); err == nil {
				//验证加密后的密码是否正确
				if passwd = crypt.Sha256(Salt + "." + passwd); login.Password != passwd {
					err = errors.New("username or password error")
				} else {
					salt := crypt.Ulid()
					token = w.generateToken(login.Username, salt, passwd, 1)
					w.secret.Add(salt, passwd)
					w.salt.Remove(login.Username)
				}
			}
		} else {
			err = errors.New(http.StatusText(http.StatusUnauthorized))
		}
	}
	ctx.JSON(http.StatusOK, api.Replay(ctx, err, token))
}

func (w *Server) DoLogin(ctx *gin.Context) {
	var doLogin dto.DoLogin
	var Salt string
	var err error
	if err = ctx.ShouldBindJSON(&doLogin); err == nil {
		if service.GetDb().CheckUserName(doLogin.Username) {
			//Salt = crypt.RandVKey()
			w.salt.Add(doLogin.Username, Salt)
		}
	}
	ctx.JSON(http.StatusOK, api.Replay(ctx, err, Salt))
}
