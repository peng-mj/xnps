package proxy

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"net"
	"net/http"
	"sort"
	"sync"
	"xnps/database/models"
	"xnps/lib/common"
	"xnps/lib/conn"
	"xnps/netManager/bridge"
)

type Service interface {
	Start() error
	Close() error
}

type NetBridge interface {
	SendLinkInfo(clientId int64, link *conn.Link, t *models.Tunnel) (target net.Conn, err error)
}

// BaseServer struct
type BaseServer struct {
	id           int
	bridge       NetBridge
	tunnel       *models.Tunnel
	errorContent []byte
	sync.Mutex
}

func NewBaseServer(brg *bridge.Bridge, task *models.Tunnel) *BaseServer {
	return &BaseServer{
		bridge:       brg,
		tunnel:       task,
		errorContent: nil,
		Mutex:        sync.Mutex{},
	}
}

// add the flow
func (s *BaseServer) FlowAdd(in, out int64) {
	s.Lock()
	defer s.Unlock()
	s.tunnel.Flow.ExportFlow += out
	s.tunnel.Flow.InletFlow += in
}

// write fail bytes to the connection
func (s *BaseServer) writeConnFail(c net.Conn) {
	c.Write([]byte(common.ConnectionFailBytes))
	c.Write(s.errorContent)
}

// auth check
func (s *BaseServer) auth(r *http.Request, c *conn.Conn, u, p string) error {
	if u != "" && p != "" && !common.CheckAuth(r, u, p) {
		c.Write([]byte(common.UnauthorizedBytes))
		c.Close()
		return errors.New("401 Unauthorized")
	}
	return nil
}

// check flow limit of the client ,and decrease the allow num of client
func (s *BaseServer) CheckFlowAndConnNum(client *models.Client) error {
	//TODO:这里检查流量限制

	if !client.GetConn() {
		return errors.New("Connections exceed the current client limit")
	}
	return nil
}

func in(target string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	if index < len(strArray) && strArray[index] == target {
		return true
	}
	return false
}

// create a new connection and start bytes copying
func (s *BaseServer) DealClient(c *conn.Conn, client *models.Client, addr string,
	rb []byte, tp string, callBack func(), localProxy bool, task *models.Tunnel) error {

	// TODO: 判断访问地址是否在黑名单内
	if common.IsBlackIp(c.RemoteAddr().String(), client.AccessKey) {
		c.Close()
		return nil
	}

	link := conn.NewLink(tp, addr, client.Crypt, client.Compress, c.Conn.RemoteAddr().String(), localProxy)
	if target, err := s.bridge.SendLinkInfo(client.Id, link, s.tunnel); err != nil {
		logs.Warn("get connection from client id %d  error %s", client.Id, err.Error())
		c.Close()
		return err
	} else {
		if callBack != nil {
			callBack()
		}
		conn.CopyWaitGroup(target, c.Conn, link.Crypt, link.Compress, client.Rate, true, rb, task)
	}
	return nil
}
