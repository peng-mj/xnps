package client

import (
	"bytes"
	"net"
	"strconv"
	"sync"
	"time"
	"xnps/lib/nps_mux"

	"github.com/astaxie/beego/logs"
	"github.com/xtaci/kcp-go"

	"xnps/lib/common"
	"xnps/lib/config"
	"xnps/lib/conn"
	"xnps/lib/crypt"
)

type TRPClient struct {
	svrAddr        string
	bridgeConnType string
	proxyUrl       string
	vKey           string
	p2pAddr        map[string]string
	tunnel         *nps_mux.Mux
	signal         *conn.Conn
	ticker         *time.Ticker
	cnf            *config.Config
	disconnectTime int
	once           sync.Once
}

// new client
func NewRPClient(svraddr string, vKey string, bridgeConnType string, proxyUrl string, cnf *config.Config, disconnectTime int) *TRPClient {
	return &TRPClient{
		svrAddr:        svraddr,
		p2pAddr:        make(map[string]string, 0),
		vKey:           vKey,
		bridgeConnType: bridgeConnType,
		proxyUrl:       proxyUrl,
		cnf:            cnf,
		disconnectTime: disconnectTime,
		once:           sync.Once{},
	}
}

var NowStatus int
var CloseClient bool

// start
func (s *TRPClient) Start() {
	CloseClient = false
	var err error
	var c *conn.Conn
retry:
	if CloseClient {
		return
	}
	NowStatus = 0
	c, err = NewConn(s.bridgeConnType, s.vKey, s.svrAddr, common.WORK_MAIN, s.proxyUrl)
	if err != nil {
		logs.Error("The connection server failed and will be reconnected in 5 seconds, error", err.Error())
		time.Sleep(time.Second * 5)
		goto retry
	}
	if c == nil {
		logs.Error("Error data from server, and will be reconnected in 5 seconds")
		time.Sleep(time.Second * 5)
		goto retry
	}
	logs.Info("Successful connection with server %s", s.svrAddr)
	//monitor the connection
	go s.ping()
	s.signal = c
	//start a channel connection
	go s.newChannel()
	//start health check if the it's open
	//if s.cnf != nil && len(s.cnf.Healths) > 0 {
	//	go heathCheck(s.cnf.Healths, s.signal)
	//}
	NowStatus = 1
	//msg connection, eg udp
	s.handleMain()
}

// handle main connection
func (s *TRPClient) handleMain() {
	for {
		flags, err := s.signal.ReadFlag()
		if err != nil {
			logs.Error("Accept server data error %s, end this service", err.Error())
			break
		}
		switch flags {
		case common.NEW_UDP_CONN:
			//read server udp addr and password
			if lAddr, err := s.signal.GetShortLenContent(); err != nil {
				logs.Warn(err)
				return
			} else if pwd, err := s.signal.GetShortLenContent(); err == nil {
				var localAddr string
				//The local port remains unchanged for a certain period of time
				if v, ok := s.p2pAddr[crypt.Sha1(string(pwd)+strconv.Itoa(int(time.Now().Unix()/100)))]; !ok {
					tmpConn, err := common.GetLocalUdpAddr()
					if err != nil {
						logs.Error(err)
						return
					}
					localAddr = tmpConn.LocalAddr().String()
				} else {
					localAddr = v
				}
				go s.newUdpConn(localAddr, string(lAddr), string(pwd))
			}
		}
	}
	s.Close()
}

func (s *TRPClient) newUdpConn(localAddr, rAddr string, md5Password string) {
	var localConn net.PacketConn
	var err error
	var remoteAddress string
	if remoteAddress, localConn, err = handleP2PUdp(localAddr, rAddr, md5Password, common.WORK_P2P_PROVIDER); err != nil {
		logs.Error(err)
		return
	}
	l, err := kcp.ServeConn(nil, 150, 3, localConn)
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Trace("start local p2p udp listen, local address", localConn.LocalAddr().String())
	for {
		udpTunnel, err := l.AcceptKCP()
		if err != nil {
			logs.Error(err)
			l.Close()
			return
		}
		if udpTunnel.RemoteAddr().String() == string(remoteAddress) {
			conn.SetUdpSession(udpTunnel)
			logs.Trace("successful connection with client ,address %s", udpTunnel.RemoteAddr().String())
			//read link info from remote
			conn.Accept(nps_mux.NewMux(udpTunnel, s.bridgeConnType, s.disconnectTime), func(c net.Conn) {
				go s.handleChannel(c)
			})
			break
		}
	}
}

// 创建新的隧道连接，也需要验证身份
// pmux tunnel
func (s *TRPClient) newChannel() {
	tunnel, err := NewConn(s.bridgeConnType, s.vKey, s.svrAddr, common.WORK_CHAN, s.proxyUrl)
	if err != nil {
		logs.Error("connect to ", s.svrAddr, "error:", err)
		return
	}
	s.tunnel = nps_mux.NewMux(tunnel.Conn, s.bridgeConnType, s.disconnectTime)
	//持续接收来自服务端的连接
	for {
		src, err := s.tunnel.Accept()
		if err != nil {
			logs.Warn(err)
			s.Close()
			break
		}
		go s.handleChannel(src)
	}
}

// 当新的隧道连接的时候，执行此函数
func (s *TRPClient) handleChannel(srcLink net.Conn) {
	link, err := conn.NewConn(srcLink).GetLinkInfo()
	if err != nil || link == nil {
		srcLink.Close()
		logs.Error("get connection info from server error ", err)
		return
	}
	//host for target processing
	link.Host = common.FormatAddress(link.Host)
	//if Conn type is http, read the request and log
	//logs.Info("type:", link.ConnType)
	//if link.ConnType == "http" {
	//	//先对目标网络建立连接
	//	if targetConnIo, err := net.DialTimeout(common.CONN_TCP, link.Host, link.Option.Timeout); err != nil {
	//		logs.Warn("connect to %s error %s", link.Host, err.Error())
	//		srcLink.Close()
	//	} else {
	//		srcConnIo := conn.GetConn(srcLink, link.Crypt, link.Compress, nil, false)
	//		//两个连接的底层实现
	//		go func() {
	//			common.CopyConnectionBuffer(srcConnIo, targetConnIo)
	//			srcConnIo.Close()
	//			targetConnIo.Close()
	//		}()
	//		for {
	//			if r, err := http.ReadRequest(bufio.NewReader(srcConnIo)); err != nil {
	//				srcConnIo.Close()
	//				targetConnIo.Close()
	//				break
	//			} else {
	//				logs.Trace("http request, method %s, host %s, url %s, remote address %s", r.Method, r.Host, r.URL.Path, r.RemoteAddr)
	//				r.Write(targetConnIo)
	//			}
	//		}
	//	}
	//	return
	//}
	if link.ConnType == "udp5" {
		logs.Trace("new %s connection with the goal of %s, remote address:%s", link.ConnType, link.Host, link.RemoteAddr)
		s.handleUdp(srcLink)
	}
	//connect to target if conn type is tcp or udp
	if targetConn, err := net.DialTimeout(link.ConnType, link.Host, link.Option.Timeout); err != nil {
		logs.Warn("connect to %s error %s", link.Host, err.Error())
		srcLink.Close()
	} else {
		logs.Trace("new %s connection with the goal of %s, remote address:%s", link.ConnType, link.Host, link.RemoteAddr)
		conn.CopyWaitGroup(srcLink, targetConn, link.Crypt, link.Compress, nil, false, nil, nil)
	}
}

func (s *TRPClient) handleUdp(serverConn net.Conn) {
	// bind a local udp port
	local, err := net.ListenUDP("udp", nil)
	defer serverConn.Close()
	if err != nil {
		logs.Error("bind local udp port error ", err.Error())
		return
	}
	defer local.Close()
	go func() {
		defer serverConn.Close()
		b := common.BufPoolUdp.Get().([]byte)
		defer common.BufPoolUdp.Put(b)
		for {
			n, raddr, err := local.ReadFrom(b)
			if err != nil {
				logs.Error("read data from remote server error", err.Error())
			}
			buf := bytes.Buffer{}
			dgram := common.NewUDPDatagram(common.NewUDPHeader(0, 0, common.ToSocksAddr(raddr)), b[:n])
			dgram.Write(&buf)
			b, err := conn.GetLengthAndBytes(buf.Bytes())
			if err != nil {
				logs.Warn("get len bytes error", err.Error())
				continue
			}
			if _, err := serverConn.Write(b); err != nil {
				logs.Error("write data to remote  error", err.Error())
				return
			}
		}
	}()
	b := common.BufPoolUdp.Get().([]byte)
	defer common.BufPoolUdp.Put(b)
	for {
		n, err := serverConn.Read(b)
		if err != nil {
			logs.Error("read udp data from server error ", err.Error())
			return
		}

		udpData, err := common.ReadUDPDatagram(bytes.NewReader(b[:n]))
		if err != nil {
			logs.Error("unpack data error", err.Error())
			return
		}
		raddr, err := net.ResolveUDPAddr("udp", udpData.Header.Addr.String())
		if err != nil {
			logs.Error("build remote addr err", err.Error())
			continue // drop silently
		}
		_, err = local.WriteTo(udpData.Data, raddr)
		if err != nil {
			logs.Error("write data to remote ", raddr.String(), "error", err.Error())
			return
		}
	}
}

// Whether the monitor channel is closed
func (s *TRPClient) ping() {
	s.ticker = time.NewTicker(time.Second * 5)
loop:
	for {
		select {
		case <-s.ticker.C:
			if s.tunnel != nil && s.tunnel.IsClose {
				s.Close()
				break loop
			}
		}
	}
}

func (s *TRPClient) Close() {
	s.once.Do(s.closing)
}

func (s *TRPClient) closing() {
	CloseClient = true
	NowStatus = 0
	if s.tunnel != nil {
		_ = s.tunnel.Close()
	}
	if s.signal != nil {
		_ = s.signal.Close()
	}
	if s.ticker != nil {
		s.ticker.Stop()
	}
}
