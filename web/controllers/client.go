package controllers

import (
	"github.com/astaxie/beego"
	"xnps/lib/common"
	"xnps/lib/database"
	"xnps/lib/database/models"
	"xnps/lib/rate"
	"xnps/server"
)

// ClientController 客户端管理相关API
type ClientController struct {
	BaseController
}

func (s *ClientController) List() {
	if s.Ctx.Request.Method == "GET" {
		s.Data["menu"] = "client"
		s.SetInfo("client")
		s.display("client/list")
		return
	}
	start, length := s.GetAjaxParams()
	clientIdSession := s.GetSession("clientId")
	var clientId int
	if clientIdSession == nil {
		clientId = 0
	} else {
		clientId = clientIdSession.(int)
	}
	list, cnt := server.GetClientList(start, length, s.getEscapeString("search"), s.getEscapeString("sort"), s.getEscapeString("order"), clientId)
	cmd := make(map[string]interface{})
	ip := s.Ctx.Request.Host
	cmd["ip"] = common.GetIpByAddr(ip)
	cmd["bridgeType"] = beego.AppConfig.String("bridge_type")
	cmd["bridgePort"] = server.Bridge.TunnelPort
	s.AjaxTable(list, cnt, cnt, cmd)
}

// 添加客户端
func (s *ClientController) Add() {
	if s.Ctx.Request.Method == "GET" {
		s.Data["menu"] = "client"
		s.SetInfo("add client")
		s.display()
	} else {
		//id := database.GetDb().JsonDb.GetClientId()
		t := &models.Client{
			VerifyKey: s.getEscapeString("vkey"),
			//Id:        id,
			Valid: true,
			Name:  s.getEscapeString("remark"),
			//HttpUser:   s.getEscapeString("u"),
			//HttpPasswd: s.getEscapeString("p"),
			//Cnf: &models.Config{
			//	User:     s.getEscapeString("u"),
			//	Passwd:   s.getEscapeString("p"),
			//	Compress: common.GetBoolByStr(s.getEscapeString("compress")),
			//	Crypt:    s.GetBoolNoErr("crypt"),
			//},
			//AllowUseConfigFile: s.GetBoolNoErr("config_conn_allow"),
			AllowUseConfigFile: false, //不允许用户使用配置文件登录
			RateLimit:          int(s.GetIntNoErr("rate_limit")),
			MaxConn:            int(s.GetIntNoErr("max_conn")),
			Compress:           common.GetBoolByStr(s.getEscapeString("compress")),
			Crypt:              s.GetBoolNoErr("crypt"),

			HttpUser:     s.getEscapeString("web_username"),
			HttpPasswd:   s.getEscapeString("web_password"),
			MaxTunnelNum: int(s.GetIntNoErr("max_tunnel")),
			Flow: &models.Flow{
				ExportFlow: 0,
				InletFlow:  0,
				FlowLimit:  int64(s.GetIntNoErr("flow_limit")),
			},
			//BlackIpList: RemoveRepeatedElement(strings.Split(s.getEscapeString("blackiplist"), "\r\n")),
		}
		if err := database.GetDb().NewClient(t); err != nil {
			s.AjaxErr(err.Error())
		}
		s.AjaxOkWithId("add success", t.Id)
	}
}
func (s *ClientController) GetClient() {
	if s.Ctx.Request.Method == "POST" {
		id := s.GetIntNoErr("id")
		data := make(map[string]interface{})
		if c, err := database.GetDb().GetClientById(id); err != nil {
			data["code"] = 0
		} else {
			data["code"] = 1
			data["data"] = c
		}
		s.Data["json"] = data
		s.ServeJSON()
	}
}

// 修改客户端
func (s *ClientController) Edit() {
	id := s.GetIntNoErr("id")
	if s.Ctx.Request.Method == "GET" {
		s.Data["menu"] = "client"
		if c, err := database.GetDb().GetClientById(id); err != nil {
			s.error()
		} else {
			s.Data["c"] = c
			//s.Data["BlackIpList"] = strings.Join(c.BlackIpList, "\r\n")
			//s.Data["BlackIpList"] = strings.Join(c.BlackIpList, "\r\n")
		}
		s.SetInfo("edit client")
		s.display()
	} else {
		if c, err := database.GetDb().GetClientById(id); err != nil {
			s.error()
			s.AjaxErr("client ID not found")
			return
		} else {
			if s.getEscapeString("web_username") != "" {
				if s.getEscapeString("web_username") == beego.AppConfig.String("web_username") || !database.GetDb().VerifyUserName(s.getEscapeString("web_username"), c.Id) {
					s.AjaxErr("web login username duplicate, please reset")
					return
				}
			}
			if s.GetSession("isAdmin").(bool) {
				if !database.GetDb().VerifyVkey(s.getEscapeString("vkey"), c.Id) {
					s.AjaxErr("Vkey duplicate, please reset")
					return
				}
				c.VerifyKey = s.getEscapeString("vkey")
				c.Flow.FlowLimit = int64(s.GetIntNoErr("flow_limit"))
				c.RateLimit = int(s.GetIntNoErr("rate_limit"))
				c.MaxConn = int(int32(s.GetIntNoErr("max_conn")))
				c.MaxTunnelNum = int(s.GetIntNoErr("max_tunnel"))
			}
			c.Name = s.getEscapeString("remark")
			c.HttpUser = s.getEscapeString("u")
			c.HttpPasswd = s.getEscapeString("p")
			c.Compress = common.GetBoolByStr(s.getEscapeString("compress"))
			c.Crypt = s.GetBoolNoErr("crypt")
			b, err := beego.AppConfig.Bool("allow_user_change_username")
			if s.GetSession("isAdmin").(bool) || (err == nil && b) {
				c.HttpUser = s.getEscapeString("web_username")
			}
			c.HttpPasswd = s.getEscapeString("web_password")
			c.AllowUseConfigFile = s.GetBoolNoErr("config_conn_allow")
			if c.Rate != nil {
				c.Rate.Stop()
			}
			if c.RateLimit > 0 {
				c.Rate = rate.NewRate(int64(c.RateLimit * 1024))
				c.Rate.Start()
			} else {
				c.Rate = rate.NewRate(int64(2 << 23))
				c.Rate.Start()
			}
			//TODO:黑名单管理需要重构
			//c.BlackIpList = RemoveRepeatedElement(strings.Split(s.getEscapeString("blackiplist"), "\r\n"))
			database.GetDb().UpdateClientById(c, c.Id)
		}
		s.AjaxOk("save success")
	}
}

func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

// 更改状态
func (s *ClientController) ChangeStatus() {
	id := s.GetIntNoErr("id")
	if client, err := database.GetDb().GetClientById(id); err == nil {
		client.Valid = s.GetBoolNoErr("status")
		if client.Valid == false {
			server.DelClientConnect(client.Id)
		}
		s.AjaxOk("modified success")
	}
	s.AjaxErr("modified fail")
}

// 删除客户端
func (s *ClientController) Del() {
	id := s.GetIntNoErr("id")
	if err := database.GetDb().DelClient(id); err != nil {
		s.AjaxErr("delete error")
	}
	//server.DelTunnelAndHostByClientId(id, false)
	server.DelClientConnect(id)
	s.AjaxOk("delete success")
}
