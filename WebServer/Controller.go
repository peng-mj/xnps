package WebServer

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4/middleware"
	"xnps/WebServer/WebApi"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func (w *WebServer) InitSystem(pre string) {
	r := w.e.Group(pre)
	r = w.e.Group("/api")
	//conf := echojwt.Config{SigningKey: []byte("secret")}

	//r.Use(w.keyAuthMiddleware)
	//r.Use(w.jwtMiddleware, echojwt.WithConfig(conf))
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.POST("/system/addConfig", WebApi.AddSysConfig)
	r.POST("/system/editConfig", WebApi.EditSYsConfig)
	r.POST("/system/Config", WebApi.GetSysConfig)
	////r.POST("/roastInfo/devId", WebApi.GetAllDevInfo)
	//r.POST("/get/tempData/devId", WebApi.GetTempInfo)

}

func (w *WebServer) ApiRouterRules(pre string) {
	r := w.e.Group(pre)
	r = w.e.Group("/api")
	//conf := echojwt.Config{SigningKey: []byte("secret")}

	//r.Use(w.keyAuthMiddleware)
	r.Use(w.jwtMiddleware)
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())

	r.POST("/roastInfo/devId", WebApi.GetAllGroup)
	r.POST("/get/tempData/devId", WebApi.GetAllClient)
	r.POST("/get/tempData/devId", WebApi.GetAllFirewall)
	r.POST("/get/tempData/devId", WebApi.GetAllBlockList)

}
