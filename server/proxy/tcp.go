package proxy

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"xnps/bridge"
	"xnps/database/models"
	"xnps/lib/common"
	"xnps/lib/conn"
	"xnps/server/connection"
)

type TunnelModeServer struct {
	BaseServer
	process  process
	listener net.Listener
}

// tcp|http|host
func NewTunnelModeServer(process process, bridge NetBridge, task *models.Tunnel) *TunnelModeServer {
	s := new(TunnelModeServer)
	s.bridge = bridge
	s.process = process
	s.tunnel = task
	return s
}

// START
func (s *TunnelModeServer) Start() error {
	return conn.NewTcpListenerAndProcess(s.tunnel.ServerIp+":"+strconv.Itoa(int(s.tunnel.ServerPort)), func(c net.Conn) {
		if err := s.CheckFlowAndConnNum(s.tunnel.Client); err != nil {
			logs.Warn("%s:%d client:%d tcp connect error:%s", c.RemoteAddr(), s.tunnel.ServerPort, s.tunnel.Client.Id, err.Error())
			c.Close()
			return
		}
		logs.Trace("%s:%d new tcp connection,client:%d", c.RemoteAddr(), s.tunnel.ServerPort, s.tunnel.Client.Id)
		s.process(conn.NewConn(c), s)
		s.tunnel.Client.AddConn()
	}, &s.listener)
}

// close
func (s *TunnelModeServer) Close() error {
	return s.listener.Close()
}

// web管理方式
type WebServer struct {
	BaseServer
}

// 开始
func (s *WebServer) Start() error {
	p, _ := beego.AppConfig.Int("web_port")
	if p == 0 {
		stop := make(chan struct{})
		<-stop
	}
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.SetStaticPath(beego.AppConfig.String("web_base_url")+"/static", filepath.Join(common.GetRunPath(), "web", "static"))
	beego.SetViewsPath(filepath.Join(common.GetRunPath(), "web", "views"))
	err := errors.New("Web management startup failure ")
	var l net.Listener
	if l, err = connection.GetWebManagerListener(); err == nil {
		beego.InitBeforeHTTPRun()
		if beego.AppConfig.String("web_open_ssl") == "true" {
			keyPath := beego.AppConfig.String("web_key_file")
			certPath := beego.AppConfig.String("web_cert_file")
			err = http.ServeTLS(l, beego.BeeApp.Handlers, certPath, keyPath)
		} else {
			err = http.Serve(l, beego.BeeApp.Handlers)
		}
	} else {
		logs.Error(err)
	}
	return err
}

func (s *WebServer) Close() error {
	return nil
}

// new
func NewWebServer(bridge *bridge.Bridge) *WebServer {
	s := new(WebServer)
	s.bridge = bridge
	return s
}

type process func(c *conn.Conn, s *TunnelModeServer) error

// tcp proxy
func ProcessTunnel(c *conn.Conn, s *TunnelModeServer) error {

	targetAddr, err := s.tunnel.Target.GetRandomTarget()
	if err != nil {
		c.Close()
		logs.Warn("tcp port %d ,client id %d,tunnel id %d connect error %s", s.tunnel.ServerPort, s.tunnel.ClientId, s.tunnel.Id, err.Error())
		return err
	}
	return s.DealClient(c, s.tunnel.Client, targetAddr, nil, common.CONN_TCP, nil, s.tunnel.Target.LocalProxy, s.tunnel)
}
