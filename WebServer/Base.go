package WebServer

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"sync"
	"xnps/lib/cache"
)

type WebServer struct {
	e        *echo.Echo
	g        *gin.Engine
	salt     *cache.Cache
	JWTToken *jwt.Token
	secret   *cache.Cache
	//tokenMan TokenManager
}

////go:embed web/static/*
//var staticFiles embed.FS

// 初始化系统时，运行的web段配置
func InitSystem(wg *sync.WaitGroup, url string) {

}

func (w *WebServer) Start(url string) {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Done()
	w.g = gin.Default()
	w.salt = cache.New(20)
	w.secret = cache.New(20) //默认最大密钥存储的文件数量
	w.JWTToken = jwt.New(jwt.SigningMethodHS256)
	//w.e = echo.New()

	//w.e.HTTPErrorHandler = w.ErrorHandler
	//装载静态文件
	//TODO:后续使用go:embed打包
	w.g.Static("/", "web/static/")
	w.g.POST("/doLogin", w.DoLogin)
	w.g.POST("/login", w.Login)
	w.ApiRouterRules("/api")
	w.e.Logger.Fatal(w.e.Start(url))
}
