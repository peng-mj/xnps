package controllers

import (
	"xnps/lib/database"
	"xnps/lib/database/models"
	"xnps/server"
	"xnps/server/tool"

	"github.com/astaxie/beego"
)

type IndexController struct {
	BaseController
}

func (s *IndexController) Index() {
	s.Data["web_base_url"] = beego.AppConfig.String("web_base_url")
	s.Data["data"] = server.GetDashboardData()
	s.SetInfo("dashboard")
	s.display("index/index")
}
func (s *IndexController) Help() {
	s.SetInfo("about")
	s.display("index/help")
}

func (s *IndexController) Tcp() {
	s.SetInfo("tcp")
	s.SetType("tcp")
	s.display("index/list")
}

func (s *IndexController) Udp() {
	s.SetInfo("udp")
	s.SetType("udp")
	s.display("index/list")
}

func (s *IndexController) Socks5() {
	s.SetInfo("socks5")
	s.SetType("socks5")
	s.display("index/list")
}

func (s *IndexController) Http() {
	s.SetInfo("http proxy")
	s.SetType("httpProxy")
	s.display("index/list")
}
func (s *IndexController) File() {
	s.SetInfo("file server")
	s.SetType("file")
	s.display("index/list")
}

func (s *IndexController) Secret() {
	s.SetInfo("secret")
	s.SetType("secret")
	s.display("index/list")
}
func (s *IndexController) P2p() {
	s.SetInfo("p2p")
	s.SetType("p2p")
	s.display("index/list")
}

func (s *IndexController) Host() {
	s.SetInfo("host")
	s.SetType("hostServer")
	s.display("index/list")
}

func (s *IndexController) All() {
	s.Data["menu"] = "client"
	clientId := s.getEscapeString("client_id")
	s.Data["client_id"] = clientId
	s.SetInfo("client id:" + clientId)
	s.display("index/list")
}

func (s *IndexController) GetTunnel() {
	start, length := s.GetAjaxParams()
	taskType := s.getEscapeString("type")
	clientId := s.GetIntNoErr("client_id")
	list, cnt := server.GetTunnel(start, length, taskType, clientId, s.getEscapeString("search"))
	s.AjaxTable(list, cnt, cnt, nil)
}

func (s *IndexController) Add() {
	if s.Ctx.Request.Method == "GET" {
		s.Data["type"] = s.getEscapeString("type")
		s.Data["client_id"] = s.getEscapeString("client_id")
		s.SetInfo("add tunnel")
		s.display()
	} else {
		id := database.GetDb().JsonDb.GetTaskId()
		t := &models.Tunnel{
			Port:      s.GetIntNoErr("port"),
			ServerIp:  s.getEscapeString("server_ip"),
			Mode:      s.getEscapeString("type"),
			Target:    &models.Target{TargetStr: s.getEscapeString("target"), LocalProxy: s.GetBoolNoErr("local_proxy")},
			Id:        id,
			Status:    true,
			Remark:    s.getEscapeString("remark"),
			Password:  s.getEscapeString("password"),
			LocalPath: s.getEscapeString("local_path"),
			StripPre:  s.getEscapeString("strip_pre"),
			Flow:      &models.Flow{},
		}

		if t.Port <= 0 {
			t.Port = tool.GenerateServerPort(t.Mode)
		}

		if !tool.TestServerPort(t.Port, t.Mode) {
			s.AjaxErr("The port cannot be opened because it may has been occupied or is no longer allowed.")
		}
		var err error
		if t.Client, err = database.GetDb().GetClient(s.GetIntNoErr("client_id")); err != nil {
			s.AjaxErr(err.Error())
		}
		if t.Client.MaxTunnelNum != 0 && t.Client.GetTunnelNum() >= t.Client.MaxTunnelNum {
			s.AjaxErr("The number of tunnels exceeds the limit")
		}
		if err := database.GetDb().NewTask(t); err != nil {
			s.AjaxErr(err.Error())
		}
		if err := server.AddTask(t); err != nil {
			s.AjaxErr(err.Error())
		} else {
			s.AjaxOkWithId("add success", id)
		}
	}
}
func (s *IndexController) GetOneTunnel() {
	id := s.GetIntNoErr("id")
	data := make(map[string]interface{})
	if t, err := database.GetDb().GetTaskById(id); err != nil {
		data["code"] = 0
	} else {
		data["code"] = 1
		data["data"] = t
	}
	s.Data["json"] = data
	s.ServeJSON()
}
func (s *IndexController) Edit() {
	id := s.GetIntNoErr("id")
	if s.Ctx.Request.Method == "GET" {
		if t, err := database.GetDb().GetTaskById(id); err != nil {
			s.error()
		} else {
			s.Data["t"] = t
		}
		s.SetInfo("edit tunnel")
		s.display()
	} else {
		if t, err := database.GetDb().GetTaskById(id); err != nil {
			s.error()
		} else {
			if client, err := database.GetDb().GetClient(s.GetIntNoErr("client_id")); err != nil {
				s.AjaxErr("modified error,the client is not exist")
				return
			} else {
				t.Client = client
			}
			if s.GetIntNoErr("port") != t.Port {
				t.Port = s.GetIntNoErr("port")

				if t.Port <= 0 {
					t.Port = tool.GenerateServerPort(t.Mode)
				}

				if !tool.TestServerPort(s.GetIntNoErr("port"), t.Mode) {
					s.AjaxErr("The port cannot be opened because it may has been occupied or is no longer allowed.")
					return
				}
			}
			t.ServerIp = s.getEscapeString("server_ip")
			t.Mode = s.getEscapeString("type")
			t.Target = &models.Target{TargetStr: s.getEscapeString("target")}
			t.Password = s.getEscapeString("password")
			t.Id = id
			t.LocalPath = s.getEscapeString("local_path")
			t.StripPre = s.getEscapeString("strip_pre")
			t.Remark = s.getEscapeString("remark")
			t.Target.LocalProxy = s.GetBoolNoErr("local_proxy")
			database.GetDb().UpdateTask(t)
			server.StopServer(t.Id)
			server.StartTask(t.Id)
		}
		s.AjaxOk("modified success")
	}
}

func (s *IndexController) Stop() {
	id := s.GetIntNoErr("id")
	if err := server.StopServer(id); err != nil {
		s.AjaxErr("stop error")
	}
	s.AjaxOk("stop success")
}

func (s *IndexController) Del() {
	id := s.GetIntNoErr("id")
	if err := server.DelTask(id); err != nil {
		s.AjaxErr("delete error")
	}
	s.AjaxOk("delete success")
}

func (s *IndexController) Start() {
	id := s.GetIntNoErr("id")
	if err := server.StartTask(id); err != nil {
		s.AjaxErr("start error")
	}
	s.AjaxOk("start success")
}
