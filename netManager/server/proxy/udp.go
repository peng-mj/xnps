package proxy

import (
	"github.com/astaxie/beego/logs"
	"io"
	"net"
	"strings"
	"sync"
	"time"
	"xnps/database/models"
	"xnps/lib/common"
	"xnps/lib/conn"
	"xnps/netManager/bridge"
)

type UdpModeServer struct {
	BaseServer
	addrMap  sync.Map
	listener *net.UDPConn
}

func NewUdpModeServer(bridge *bridge.Bridge, task *models.Tunnel) *UdpModeServer {
	s := new(UdpModeServer)
	s.bridge = bridge
	s.tunnel = task
	return s
}

// 开始
func (s *UdpModeServer) Start() error {
	var err error
	if s.tunnel.ServerIp == "" {
		s.tunnel.ServerIp = "0.0.0.0"
	}
	s.listener, err = net.ListenUDP("udp", &net.UDPAddr{net.ParseIP(s.tunnel.ServerIp), int(s.tunnel.ServerPort), ""})
	if err != nil {
		return err
	}
	for {
		buf := common.BufPoolUdp.Get().([]byte)
		n, addr, err := s.listener.ReadFromUDP(buf)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
			continue
		}

		// 判断访问地址是否在黑名单内
		//黑名单管理需重构
		if common.IsBlackIp(addr.String(), s.tunnel.Client.AccessKey) {
			break
		}

		logs.Trace("New udp connection,client %d,remote address %s", s.tunnel.Client.Id, addr)
		go s.process(addr, buf[:n])
	}
	return nil
}

func (s *UdpModeServer) process(addr *net.UDPAddr, data []byte) {
	if v, ok := s.addrMap.Load(addr.String()); ok {
		clientConn, ok := v.(io.ReadWriteCloser)
		if ok {
			_, err := clientConn.Write(data)
			if err != nil {
				logs.Warn(err)
				return
			}
			//流量记录
			//s.tunnel.Client.Flow.Add(int64(len(data)), int64(len(data)))
		}
	} else {
		if err := s.CheckFlowAndConnNum(s.tunnel.Client); err != nil {
			logs.Warn("client id %d, tunnel id %d,error %s, when udp connection", s.tunnel.Client.Id, s.tunnel.Id, err.Error())
			return
		}
		defer s.tunnel.Client.AddConn()
		link := conn.NewLink(common.CONN_UDP, s.tunnel.Target.TargetStr, s.tunnel.Client.Crypt, s.tunnel.Client.Compress, addr.String(), s.tunnel.Target.LocalProxy)
		if clientConn, err := s.bridge.SendLinkInfo(s.tunnel.Client.Id, link, s.tunnel); err != nil {
			return
		} else {
			target := conn.GetConn(clientConn, s.tunnel.Client.Crypt, s.tunnel.Client.Compress, nil, true)
			s.addrMap.Store(addr.String(), target)
			defer target.Close()

			_, err := target.Write(data)
			if err != nil {
				logs.Warn(err)
				return
			}

			buf := common.BufPoolUdp.Get().([]byte)
			defer common.BufPoolUdp.Put(buf)
			//流量记录
			//s.tunnel.Client.Flow.Add(int64(len(data)), int64(len(data)))
			for {
				clientConn.SetReadDeadline(time.Now().Add(time.Minute * 10))
				if n, err := target.Read(buf); err != nil {
					s.addrMap.Delete(addr.String())
					logs.Warn(err)
					return
				} else {
					_, err := s.listener.WriteTo(buf[:n], addr)
					if err != nil {
						logs.Warn(err)
						return
					}
					//流量记录

					//s.tunnel.Client.Flow.Add(int64(n), int64(n))
				}
				if err := s.CheckFlowAndConnNum(s.tunnel.Client); err != nil {
					logs.Warn("client id %d, tunnel id %d,error %s, when udp connection", s.tunnel.Client.Id, s.tunnel.Id, err.Error())
					return
				}
			}
		}
	}
}

func (s *UdpModeServer) Close() error {
	return s.listener.Close()
}
