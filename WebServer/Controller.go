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

func (w *WebServer) ApiRouterRules(pre string) {
	r := w.e.Group(pre)
	r = w.e.Group("/api")
	r.Use(w.jwtMiddleware)
	//分组管理
	r.POST("/group/get/all", WebApi.GetAllGroup)
	r.POST("/group/get/condition", WebApi.GetGroupByCondition)
	r.POST("/group/add/condition", WebApi.AddGroup)
	r.POST("/group/update/one", WebApi.EditGroup)
	r.POST("/group/delete/one", WebApi.DelGroup)
	//终端相关
	r.POST("/client/get/all", WebApi.GetAllClient)
	r.POST("/client/get/condition", WebApi.GetGroupByCondition)
	r.POST("/client/add/condition", WebApi.AddClient)
	r.POST("/client/update/one", WebApi.EditClient)
	r.POST("/client/delete/one", WebApi.DelClient)
	//隧道相关
	r.POST("/tunnel/get/all", WebApi.GetAllTunnel)
	r.POST("/tunnel/get/condition", WebApi.GetTunnelByCondition)
	r.POST("/tunnel/add/condition", WebApi.AddTunnel)
	r.POST("/tunnel/update/one", WebApi.EditTunnel)
	r.POST("/tunnel/delete/one", WebApi.DelTunnel)
	//防火墙
	r.POST("/firewall/get/all", WebApi.GetAllFirewall)
	r.POST("/firewall/get/condition", WebApi.GetGroupByFirewall) //这个需要弄清楚
	r.POST("/firewall/add/condition", WebApi.AddFirewall)
	r.POST("/firewall/update/one", WebApi.EditFirewall)
	r.POST("/firewall/delete/one", WebApi.DelFirewall)
	//黑名单相关
	r.POST("/block/get/all", WebApi.GetAllBlockList)
	r.POST("/block/get/condition", WebApi.GetBlockListByCondition)
	r.POST("/block/add/condition", WebApi.AddFBlockList)
	r.POST("/block/update/one", WebApi.EditBlockList)
	r.POST("/block/delete/one", WebApi.DelBlockList)

}
