package WebServer

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"sync"
	"xnps/WebServer/WebObj"
)

type WebServer struct {
	e    *echo.Echo
	salt *WebObj.Salt
}

////go:embed web/static*
//var staticFiles embed.FS

func (w *WebServer) Start(url string) {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Done()
	w.salt = WebObj.NewSaltMap(20)
	w.e = echo.New()
	w.e.HTTPErrorHandler = w.ErrorHandler
	//装载静态文件
	w.e.Static("/", "web/static/")
	//w.e.Static("/", staticFiles.ReadDir())
	w.e.Use(middleware.BodyLimit("2M"))
	w.e.Use(middleware.CORS()) //NOTE:如果跨域，需要特别注意
	w.e.POST("/doLogin", w.DoLogin)
	w.e.POST("/login", w.Login)
	w.ApiRouterRules("/api")
	//r := w.e.Group("/api")

	//r.Use(w.keyAuthMiddleware)
	//r.POST("/devInfo/all", WebApi.GetAllDevInfo)
	//r.POST("/roastInfo/devId", WebApi.GetRoastInfoById)
	////r.POST("/roastInfo/devId", WebApi.GetAllDevInfo)
	//r.POST("/get/tempData/devId", WebApi.GetTempInfo)

	w.e.Logger.Fatal(w.e.Start(url))
}

//func (w *WebServer) LoginMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
//
//	return func(c echo.Context) error {
//
//		//获取请求中的密钥
//		Salt := c.Request().Header.Get("Salt")
//		user := c.Request().Header.Get("username")
//
//		// 验证密钥是否正确
//		//if err != nil || key != token {
//		//Utils.Log("token 验证失败")
//		//re := WebObj.ReData{Data: ReCode.ErTokenLapsed, Status: "ERROR"}
//		//request, _ := json.Marshal(re)
//		//return c.HTML(http.StatusAccepted, string(request))
//		//}
//		//Utils.Log("密钥验证通过")
//		// 密钥验证通过，继续传递请求给下一个处理程序
//		return next(c)
//	}
//}

func (w *WebServer) ErrorHandler(err error, c echo.Context) {
	if he, ok := err.(*echo.HTTPError); ok {
		switch he.Code {
		case 404:
			if err := c.File("static/404.html"); err != nil {
				c.Logger().Error(err)
			}
		}
	}
}
