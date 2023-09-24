package WebServer

import (
	"github.com/golang-jwt/jwt"
	"xnps/WebServer/WebApi"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

//func (w *WebServer) InitSystem(pre string) {
//	r := w.e.Group(pre)
//	r = w.e.Group("/api")
//	//conf := echojwt.Config{SigningKey: []byte("secret")}
//
//	//r.Use(w.keyAuthMiddleware)
//	//r.Use(w.jwtMiddleware, echojwt.WithConfig(conf))
//	r.Use(middleware.Logger())
//	r.Use(middleware.Recover())
//	r.POST("/system/addConfig", WebApi.AddSysConfig)
//	r.POST("/system/editConfig", WebApi.EditSYsConfig)
//	r.POST("/system/Config", WebApi.GetSysConfig)
//	////r.POST("/roastInfo/devId", WebApi.GetAllDevInfo)
//	//r.POST("/get/tempData/devId", WebApi.GetTempInfo)
//
//}

func (w *WebServer) ApiRouterRules(pre string) {
	r := w.e.Group(pre)
	r = w.e.Group("/api")
	r.Use(w.jwtMiddleware)
	//r.Use(middleware.Logger())
	//r.Use(middleware.Recover())

	r.POST("/group/get/all", WebApi.GetAllGroup)
	r.POST("/group/get/condition", WebApi.GetGroupByCondition)
	r.POST("/group/add/condition", WebApi.AddGroup)
	r.POST("/group/update/one", WebApi.EditGroup)
	r.POST("/group/delete/one", WebApi.DelGroup)

	r.POST("/client/get/all", WebApi.GetAllClient)
	r.POST("/client/get/condition", WebApi.GetGroupByCondition)
	r.POST("/client/add/condition", WebApi.AddClient)
	r.POST("/client/update/one", WebApi.EditClient)
	r.POST("/client/delete/one", WebApi.DelClient)

	r.POST("/tunnel/get/all", WebApi.GetAllTunnel)
	r.POST("/tunnel/get/condition", WebApi.GetTunnelByCondition)
	r.POST("/tunnel/add/condition", WebApi.AddTunnel)
	r.POST("/tunnel/update/one", WebApi.EditTunnel)
	r.POST("/tunnel/delete/one", WebApi.DelTunnel)

	r.POST("/firewall/get/all", WebApi.GetAllFirewall)
	r.POST("/firewall/get/condition", WebApi.GetGroupByFirewall) //这个需要弄清楚
	r.POST("/firewall/add/condition", WebApi.AddFirewall)
	r.POST("/firewall/update/one", WebApi.EditFirewall)
	r.POST("/firewall/delete/one", WebApi.DelFirewall)

	r.POST("/firewall/get/all", WebApi.GetAllBlockList)
	r.POST("/firewall/get/condition", WebApi.GetBlockListByCondition)
	r.POST("/firewall/add/condition", WebApi.AddFBlockList)
	r.POST("/firewall/update/one", WebApi.EditBlockList)
	r.POST("/firewall/delete/one", WebApi.DelBlockList)

}
