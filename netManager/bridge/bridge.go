package bridge

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
	"xnps/database/Mapper"
	"xnps/database/models"
	"xnps/lib/nps_mux"
	"xnps/netManager/server/connection"
	"xnps/netManager/server/tool"

	"xnps/lib/common"
	"xnps/lib/conn"
	"xnps/lib/crypt"
	"xnps/lib/version"
)

type Client struct {
	tunnel    *nps_mux.Mux
	signal    *conn.Conn
	file      *nps_mux.Mux
	Version   string
	retryTime int // it will be add 1 when ping not ok until to 3 will close the client
}

func NewClient(t, f *nps_mux.Mux, s *conn.Conn, vs string) *Client {
	return &Client{
		signal:  s,
		tunnel:  t,
		file:    f,
		Version: vs,
	}
}

type Bridge struct {
	TunnelPort     int //通信隧道端口
	Client         sync.Map
	Register       sync.Map
	tunnelType     string //bridge type kcp or tcp
	OpenTask       chan *models.Tunnel
	CloseTask      chan *models.Tunnel
	CloseClient    chan int64
	SecretChan     chan *conn.Secret
	ipVerify       bool
	runList        sync.Map //map[int]interface{}
	disconnectTime int
}

func NewTunnel(tunnelPort int, tunnelType string, ipVerify bool, runList sync.Map, disconnectTime int) *Bridge {
	return &Bridge{
		TunnelPort:     tunnelPort,
		tunnelType:     tunnelType,
		OpenTask:       make(chan *models.Tunnel),
		CloseTask:      make(chan *models.Tunnel),
		CloseClient:    make(chan int64),
		SecretChan:     make(chan *conn.Secret),
		ipVerify:       ipVerify,
		runList:        runList,
		disconnectTime: disconnectTime,
	}
}

func (s *Bridge) StartTunnel() error {
	go s.ping()
	if s.tunnelType == "kcp" {
		slog.Info("server start, the bridge type is %s, the bridge port is %d", s.tunnelType, s.TunnelPort)
		//return conn.NewKcpListenerAndProcess(beego.AppConfig.String("bridge_ip")+":"+beego.AppConfig.String("bridge_port"), func(c net.Conn) {
		//	s.clientProcess(conn.NewConn(c))
		//})
		return errors.New("kcp is not support. and it will remove later")
	} else {
		listener, err := connection.GetBridgeListener(s.tunnelType)
		if err != nil {
			slog.Error("GetBridgeListener ", "error", err)
			os.Exit(0)
			return err
		}
		//接收客户端连接
		conn.Accept(listener, func(c net.Conn) {
			s.clientProcess(conn.NewConn(c))
		})
	}
	return nil
}

// 验证失败，返回错误验证flag，并且关闭连接
func (s *Bridge) verifyError(c *conn.Conn) {
	c.Write([]byte(common.VERIFY_EER))
}

// 验证成功，返回错误验证flag，并且关闭连接
func (s *Bridge) verifySuccess(c *conn.Conn) {
	c.Write([]byte(common.VERIFY_SUCCESS))
}

// 客户端连接服务
func (s *Bridge) clientProcess(c *conn.Conn) {
	//read test flag
	if _, err := c.GetShortContent(3); err != nil {
		slog.Info("The client %s connect error", c.Conn.RemoteAddr(), err.Error())
		return
	}
	//version check
	if b, err := c.GetShortLenContent(); err != nil || string(b) != version.GetCoreVersion() {
		slog.Info("The client %s version does not match", c.Conn.RemoteAddr())
		c.Close()
		return
	}
	//version get
	var ver []byte
	var err error
	if ver, err = c.GetShortLenContent(); err != nil {
		slog.Info("get client %s version error", err.Error())
		c.Close()
		return
	}
	//write server version to client
	c.Write([]byte(crypt.Sha1(version.GetCoreVersion())))
	c.SetReadDeadlineBySecond(5)
	var buf []byte
	//get vKey from client
	//如为md5那么下面数字为32
	//如为sha256那么下面的数字为64
	if buf, err = c.GetShortContent(40); err != nil {
		c.Close()
		return
	}
	//TODO:客户端验证
	//verify
	//id, err := Mapper.GetDb().GetClientByAccessUser(string(buf), c.Conn.RemoteAddr().String())
	client, err := Mapper.GetDb().GetClientByAccessUser(string(buf), string(buf))
	if err != nil {
		slog.Info("Current client connection validation error, close this client:", c.Conn.RemoteAddr())
		s.verifyError(c)
		return
	} else {
		s.verifySuccess(c)
	}
	if flag, err := c.ReadFlag(); err == nil {
		s.typeDeal(flag, c, client.Id, string(ver))
	} else {
		slog.Warn("read flag", err, flag)
	}
	return
}

func (s *Bridge) DelClient(id int64) {
	if v, ok := s.Client.Load(id); ok {
		if v.(*Client).signal != nil {
			v.(*Client).signal.Close()
		}
		s.Client.Delete(id)
		if Mapper.GetDb().CheckClientValid(id) {
			return
		}
		if c, err := Mapper.GetDb().GetClientById(id); err == nil {
			s.CloseClient <- c.Id
		}
	}
}

// use different
func (s *Bridge) typeDeal(typeVal string, conn *conn.Conn, id int64, ver string) {
	valid := Mapper.GetDb().CheckClientValid(id)
	switch typeVal {
	case common.WORK_MAIN:
		if valid {
			conn.Close()
			return
		}
		tcpConn, ok := conn.Conn.(*net.TCPConn)
		if ok {
			// add tcp keep alive option for signal connection
			_ = tcpConn.SetKeepAlive(true)
			_ = tcpConn.SetKeepAlivePeriod(5 * time.Second)
		}
		//TODO:这里有待考虑
		//the vKey connect by another ,close the client of before
		if v, ok := s.Client.LoadOrStore(id, NewClient(nil, nil, conn, ver)); ok {
			if v.(*Client).signal != nil {
				v.(*Client).signal.WriteClose()
			}
			v.(*Client).signal = conn
			v.(*Client).Version = ver
		}
		//go s.GetHealthFromClient(id, conn)
		slog.Info("连接成功", "clientId %d connection succeeded, address:%s ", id, conn.Conn.RemoteAddr())

	//	TODO:隧道连接在这里
	case common.WORK_CHAN:
		muxConn := nps_mux.NewMux(conn.Conn, s.tunnelType, s.disconnectTime)
		if v, ok := s.Client.LoadOrStore(id, NewClient(muxConn, nil, nil, ver)); ok {
			v.(*Client).tunnel = muxConn
		}

	case common.WORK_CONFIG:
		client, err := Mapper.GetDb().GetClientById(id)
		if err != nil || (!valid && !client.AllowUseConfigFile) {
			conn.Close()
			return
		}
		binary.Write(conn, binary.LittleEndian, valid)
		go s.getConfig(conn, valid, client)
	case common.WORK_REGISTER:
		go s.register(conn)
	//case common.WORK_SECRET: //私密代理
	//	if passwdBytes, err := conn.GetShortContent(40); err == nil {
	//		s.SecretChan <- conn.NewSecret(string(passwdBytes), conn)
	//	} else {
	//		slog.Error("secret error, failed to match the key successfully")
	//	}
	case common.WORK_FILE:
		muxConn := nps_mux.NewMux(conn.Conn, s.tunnelType, s.disconnectTime)
		if v, ok := s.Client.LoadOrStore(id, NewClient(nil, muxConn, nil, ver)); ok {
			v.(*Client).file = muxConn
		}

	}
	conn.SetAlive(s.tunnelType)
	return
}

// register ip
func (s *Bridge) register(c *conn.Conn) {
	var hour int32
	if err := binary.Read(c, binary.LittleEndian, &hour); err == nil {
		s.Register.Store(common.GetIpByAddr(c.Conn.RemoteAddr().String()), time.Now().Add(time.Hour*time.Duration(hour)))
	}
}

func (s *Bridge) SendLinkInfo(clientId int64, link *conn.Link, t *models.Tunnel) (targetConn net.Conn, err error) {
	//if the proxy type is local
	if link.LocalProxy {
		targetConn, err = net.Dial("tcp", link.Host)
		return
	}
	if v, ok := s.Client.Load(clientId); ok {
		//If ip is restricted to do ip verification
		if s.ipVerify {
			ip := common.GetIpByAddr(link.RemoteAddr)
			if v, ok := s.Register.Load(ip); !ok {
				return nil, errors.New(fmt.Sprintf("The ip %s is not in the validation list", ip))
			} else {
				if !v.(time.Time).After(time.Now()) {
					return nil, errors.New(fmt.Sprintf("The validity of the ip %s has expired", ip))
				}
			}
		}
		var tunnel *nps_mux.Mux
		if t != nil && t.Mode == "file" { //文件代理
			tunnel = v.(*Client).file
		} else {
			tunnel = v.(*Client).tunnel
		}
		if tunnel == nil {
			err = errors.New("the client connect error")
			return
		}
		if targetConn, err = tunnel.NewConn(); err != nil {
			return
		}
		if t != nil && t.Mode == "file" {
			//TODO if t.mode is file ,not use crypt or compress
			link.Crypt = false
			link.Compress = false
			return
		}
		if _, err = conn.NewConn(targetConn).SendInfo(link, ""); err != nil {
			slog.Info("new connect error ,the targetConn %s refuse to connect", link.Host)
			return
		}
	} else {
		err = errors.New(fmt.Sprintf("the client %d is not connect", clientId))
	}
	return
}

func (s *Bridge) ping() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			arr := make([]int64, 0)
			s.Client.Range(func(key, value interface{}) bool {
				v := value.(*Client)
				if v.tunnel == nil || v.signal == nil {
					v.retryTime += 1
					if v.retryTime >= 3 {
						arr = append(arr, key.(int64))
					}
					return true
				}
				if v.tunnel.IsClose {
					arr = append(arr, key.(int64))
				}
				return true
			})
			for _, v := range arr {
				slog.Info("the client %d closed", v)
				s.DelClient(v)
			}
		}
	}
}

// TODO:去掉这个功能
// 从设备的文件启动
// get config and add task from client config
func (s *Bridge) getConfig(c *conn.Conn, isPub bool, client *models.Client) {
	var fail bool
loop:
	for {
		flag, err := c.ReadFlag()
		if err != nil {
			break
		}
		switch flag {
		case common.WORK_STATUS:
			if b, err := c.GetShortContent(64); err != nil {
				break loop
			} else {
				var str string
				id, err := Mapper.GetDb().GetClientIdByVkey(string(b))
				if err != nil {
					break loop
				}
				tunnelList, _ := Mapper.GetDb().GetTunnelListByClientId(0, id)
				//TODO:理解，为什么这么做？标记+分割号
				//应该是告诉客户端，将要创建的通道
				for i := range tunnelList {
					str += tunnelList[i].Remark + common.CONN_DATA_SEQ
				}

				_ = binary.Write(c, binary.LittleEndian, int32(len([]byte(str))))
				_ = binary.Write(c, binary.LittleEndian, []byte(str))
			}
		case common.NEW_CONF:
			var err error
			if client, err = c.GetConfigInfo(); err != nil {
				fail = true
				c.WriteAddFail()
				break loop
			} else {
				if err = Mapper.GetDb().CreateNewClient(client); err != nil {
					fail = true
					c.WriteAddFail()
					break loop
				}
				c.WriteAddOk()
				c.Write([]byte(client.AccessKey)) //这是为什么？还要向客户端写密钥？
				s.Client.Store(client.Id, NewClient(nil, nil, nil, ""))
			}
		case common.NEW_TASK:
			if tun, err := c.GetTunnelInfo(); err != nil {
				fail = true
				c.WriteAddFail()
				break loop
			} else {
				ports := common.GetPorts(tun.Ports)
				//slog.Info(ports)
				if len(ports) == 0 {
					break loop
				}
				slog.Info(tun.Target.TargetStr)
				targets := common.GetPorts(tun.Target.TargetStr)
				if len(ports) > 1 && (tun.Mode == "tcp" || tun.Mode == "udp") && (len(ports) != len(targets)) {
					fail = true
					c.WriteAddFail()
					break loop
				} else if tun.Mode == "secret" || tun.Mode == "p2p" { //限定p2p、secret才能使用
					ports = append(ports, 0)
				}
				if len(ports) == 0 {
					fail = true
					c.WriteAddFail()
					break loop
				}
				for i := 0; i < len(ports); i++ { //当端口为多个的时候，循环创建多个
					tunnel := new(models.Tunnel)
					tunnel.Mode = tun.Mode
					tunnel.ServerPort = ports[i]
					tunnel.ServerIp = tun.ServerIp
					if len(ports) == 1 {
						tunnel.Target = tun.Target
						tunnel.Remark = tun.Remark
					} else {
						tunnel.Remark = tun.Remark + "_" + strconv.Itoa(int(tunnel.ServerPort))
						tunnel.Target = new(models.Target)
						if tun.TargetAddr != "" {
							tunnel.Target.TargetStr = tun.TargetAddr + ":" + strconv.Itoa(int(targets[i]))
						} else {
							tunnel.Target.TargetStr = strconv.Itoa(int(targets[i]))
						}
					}
					//获取新的ID
					//tunnel.Id = database.GetDb().JsonDb.GetTaskId()
					tunnel.Status = true
					tunnel.Flow = &models.Flow{ClientId: tunnel.ClientId}
					//tunnel.NoStore = true
					tunnel.Client = client
					tunnel.Password = tun.Password
					tunnel.LocalPath = tun.LocalPath
					tunnel.StripPre = tun.StripPre
					//tunnel.MultiAccount = tun.MultiAccount
					//检查某客户端是否有存在的通道
					if !Mapper.GetDb().HasTunnel(client.Id, tunnel) {
						if err := Mapper.GetDb().NewTunnel(tunnel); err != nil {
							slog.Warn("Add task error ", err.Error())
							fail = true
							c.WriteAddFail()
							break loop
						}
						if b := tool.TestServerPort(tunnel.ServerPort, tunnel.Mode); !b && tun.Mode != "secret" && tun.Mode != "p2p" {
							fail = true
							c.WriteAddFail()
							break loop
						} else {
							s.OpenTask <- tunnel
						}
					}
					c.WriteAddOk()
				}
			}
		}
	}
	if fail && client != nil {
		s.DelClient(client.Id)
	}
	c.Close()
}
