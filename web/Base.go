package web

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/juju/ratelimit"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"sync"
	"time"
	"xnps/lib/crypt"
	"xnps/web/api"
	"xnps/web/dto"
	"xnps/web/service"
)

type WebServer struct {
	//e        *echo.Echo
	g *gin.Engine

	//tokenMan TokenManager
}

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func (w *WebServer) Start(url string, db *gorm.DB) {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Done()
	w.g = gin.Default()
	w.g.Use(gin.Logger(), gin.Recovery())

	//w.g.POST("/doLogin", w.DoLogin)
	//w.g.POST("/doLogin", gin.WrapF(w.DoLogin))
	w.g.GET("/ping", w.DoLogin).
		GET("/doLogin", w.DoLogin).
		POST("/login", w.Login)
	apiRouter := w.g.Group("api")
	apiRouter.Use(w.jwtMiddleware())

	w.ApiRouterRules("/Api")
	//w.e.Logger.Fatal(w.e.Start(url))
}

func (w *WebServer) ApiRouterRules(pre string) {
	r := w.g.Group(pre)
	r = w.g.Group("/Api")
	r.Use(w.jwtMiddleware())
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

}

func (w *WebServer) jwtMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 2003,
				"msg":  "",
			})
			ctx.Abort()
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if uuid, ok := claims["uuid"].(string); ok {
						if secret, ok := w.secret.GetString(uuid); ok {
							return []byte(secret), nil
						}
					}
				}
			}
			return nil, errors.New(http.StatusText(http.StatusUnauthorized))
		})
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && err == nil {
			//TODO：这个地方可有可无->Set
			ctx.Set("user", claims)
			return
		}
		ctx.Next()
	}
}

func (w *WebServer) generateToken(username, uid, passwd string, timeoutHour int) string {
	claims := w.JWTToken.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["uid"] = uid
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(timeoutHour)).Unix() // 设置Token的有效期
	tokenString, _ := w.JWTToken.SignedString([]byte(passwd))
	return tokenString
}
func (w *WebServer) RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
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
func (w *WebServer) Login(ctx *gin.Context) {
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
					salt := crypt.GetUlid()
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

func (w *WebServer) DoLogin(ctx *gin.Context) {
	var doLogin dto.DoLogin
	var Salt string
	var err error
	if err = ctx.ShouldBindJSON(&doLogin); err == nil {
		if service.GetDb().CheckUserName(doLogin.Username) {
			Salt = crypt.GenerateRandomVKey()
			w.salt.Add(doLogin.Username, Salt)
		}
	}
	ctx.JSON(http.StatusOK, api.Replay(ctx, err, Salt))
}
