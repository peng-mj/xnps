package bridge

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
	"xnps/lib/database/models"
	"xnps/lib/nps_mux"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"xnps/lib/common"
	"xnps/lib/conn"
	"xnps/lib/crypt"
	"xnps/lib/database"
	"xnps/lib/version"
	"xnps/server/connection"
	"xnps/server/tool"
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
		logs.Info("server start, the bridge type is %s, the bridge port is %d", s.tunnelType, s.TunnelPort)
		return conn.NewKcpListenerAndProcess(beego.AppConfig.String("bridge_ip")+":"+beego.AppConfig.String("bridge_port"), func(c net.Conn) {
			s.clientProcess(conn.NewConn(c))
		})
	} else {
		listener, err := connection.GetBridgeListener(s.tunnelType)
		if err != nil {
			logs.Error(err)
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

//
//// 从客户端获取运行状况
//// get health information from client
//func (s *Bridge) GetHealthFromClient(id int, c *conn.Conn) {
//	for {
//		logs.Info("health info")
//		if info, status, err := c.GetHealthInfo(); err != nil {
//			break
//		} else if !status { //the status is true , return target to the targetArr
//			file.GetDb().JsonDb.Tasks.Range(func(key, value interface{}) bool {
//				v := value.(*file.Tunnel)
//				if v.Client.Id == id && v.Mode == "tcp" && strings.Contains(v.Target.TargetStr, info) {
//					v.Lock()
//					if v.Target.TargetArr == nil || (len(v.Target.TargetArr) == 0 && len(v.HealthRemoveArr) == 0) {
//						v.Target.TargetArr = common.TrimArr(strings.Split(v.Target.TargetStr, "\n"))
//					}
//					v.Target.TargetArr = common.RemoveArrVal(v.Target.TargetArr, info)
//					if v.HealthRemoveArr == nil {
//						v.HealthRemoveArr = make([]string, 0)
//					}
//					logs.Info(info)
//					v.HealthRemoveArr = append(v.HealthRemoveArr, info)
//					v.Unlock()
//				}
//				return true
//			})
//
//		} else { //the status is false,remove target from the targetArr
//			file.GetDb().JsonDb.Tasks.Range(func(key, value interface{}) bool {
//				v := value.(*file.Tunnel)
//				if v.Client.Id == id && v.Mode == "tcp" && common.IsArrContains(v.HealthRemoveArr, info) && !common.IsArrContains(v.Target.TargetArr, info) {
//					v.Lock()
//					v.Target.TargetArr = append(v.Target.TargetArr, info)
//					v.HealthRemoveArr = common.RemoveArrVal(v.HealthRemoveArr, info)
//					v.Unlock()
//				}
//				return true
//			})
//
//		}
//	}
//	s.DelClient(id)
//}

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
		logs.Info("The client %s connect error", c.Conn.RemoteAddr(), err.Error())
		return
	}
	//version check
	if b, err := c.GetShortLenContent(); err != nil || string(b) != version.GetCoreVersion() {
		logs.Info("The client %s version does not match", c.Conn.RemoteAddr())
		c.Close()
		return
	}
	//version get
	var ver []byte
	var err error
	if ver, err = c.GetShortLenContent(); err != nil {
		logs.Info("get client %s version error", err.Error())
		c.Close()
		return
	}
	//write server version to client
	c.Write([]byte(crypt.Sha256(version.GetCoreVersion())))
	c.SetReadDeadlineBySecond(5)
	var buf []byte
	//get vKey from client
	//如为md5那么下面数字为32
	//如为sha256那么下面的数字为64
	if buf, err = c.GetShortContent(64); err != nil {
		c.Close()
		return
	}
	//TODO:客户端验证
	//verify
	id, err := database.GetDb().GetIdByVerifyKey(string(buf), c.Conn.RemoteAddr().String())
	if err != nil {
		logs.Info("Current client connection validation error, close this client:", c.Conn.RemoteAddr())
		s.verifyError(c)
		return
	} else {
		s.verifySuccess(c)
	}
	//TODO:这个flag类型，有点不懂
	if flag, err := c.ReadFlag(); err == nil {
		s.typeDeal(flag, c, id, string(ver))
	} else {
		logs.Warn(err, flag)
	}
	return
}

func (s *Bridge) DelClient(id int64) {
	if v, ok := s.Client.Load(id); ok {
		if v.(*Client).signal != nil {
			v.(*Client).signal.Close()
		}
		s.Client.Delete(id)
		if database.GetDb().IsPubClient(id) {
			return
		}
		if c, err := database.GetDb().GetClientById(id); err == nil {
			s.CloseClient <- c.Id
		}
	}
}

// use different
func (s *Bridge) typeDeal(typeVal string, c *conn.Conn, id int64, ver string) {
	isPub := database.GetDb().IsPubClient(id)
	switch typeVal {
	case common.WORK_MAIN:
		if isPub {
			c.Close()
			return
		}
		tcpConn, ok := c.Conn.(*net.TCPConn)
		if ok {
			// add tcp keep alive option for signal connection
			_ = tcpConn.SetKeepAlive(true)
			_ = tcpConn.SetKeepAlivePeriod(5 * time.Second)
		}
		//TODO:这里有待考虑
		//the vKey connect by another ,close the client of before
		if v, ok := s.Client.LoadOrStore(id, NewClient(nil, nil, c, ver)); ok {
			if v.(*Client).signal != nil {
				v.(*Client).signal.WriteClose()
			}
			v.(*Client).signal = c
			v.(*Client).Version = ver
		}
		//go s.GetHealthFromClient(id, c)
		logs.Info("clientId %d connection succeeded, address:%s ", id, c.Conn.RemoteAddr())
	//	TODO:隧道连接在这里
	case common.WORK_CHAN:
		muxConn := nps_mux.NewMux(c.Conn, s.tunnelType, s.disconnectTime)
		if v, ok := s.Client.LoadOrStore(id, NewClient(muxConn, nil, nil, ver)); ok {
			v.(*Client).tunnel = muxConn
		}

	case common.WORK_CONFIG:
		client, err := database.GetDb().GetClientById(id)
		if err != nil || (!isPub && !client.AllowUseConfigFile) {
			c.Close()
			return
		}
		binary.Write(c, binary.LittleEndian, isPub)
		go s.getConfig(c, isPub, client)
	case common.WORK_REGISTER:
		go s.register(c)
	case common.WORK_SECRET:
		if passwdBytes, err := c.GetShortContent(64); err == nil {
			s.SecretChan <- conn.NewSecret(string(passwdBytes), c)
		} else {
			logs.Error("secret error, failed to match the key successfully")
		}
	case common.WORK_FILE:
		muxConn := nps_mux.NewMux(c.Conn, s.tunnelType, s.disconnectTime)
		if v, ok := s.Client.LoadOrStore(id, NewClient(nil, muxConn, nil, ver)); ok {
			v.(*Client).file = muxConn
		}
	case common.WORK_P2P:
		//read md5 secret
		if b, err := c.GetShortContent(64); err != nil {
			logs.Error("p2p error,", err.Error())
		} else if t := database.GetDb().GetTaskByMd5Password(string(b)); t == nil {
			logs.Error("p2p error, failed to match the key successfully")
		} else {
			if v, ok := s.Client.Load(t.Client.Id); !ok {
				return
			} else {
				//向密钥对应的客户端发送与服务端udp建立连接信息，地址，密钥
				v.(*Client).signal.Write([]byte(common.NEW_UDP_CONN))
				svrAddr := beego.AppConfig.String("p2p_ip") + ":" + beego.AppConfig.String("p2p_port")
				if err != nil {
					logs.Warn("get local udp addr error")
					return
				}
				v.(*Client).signal.WriteLenContent([]byte(svrAddr))
				v.(*Client).signal.WriteLenContent(b)
				//向该请求者发送建立连接请求,服务器地址
				c.WriteLenContent([]byte(svrAddr))
			}
		}
	}
	c.SetAlive(s.tunnelType)
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
			logs.Info("new connect error ,the targetConn %s refuse to connect", link.Host)
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
				logs.Info("the client %d closed", v)
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
				id, err := database.GetDb().GetClientIdByVkey(string(b))
				if err != nil {
					break loop
				}
				tunnelList, _ := database.GetDb().GetTunnelListByClientId(0, id)
				//TODO:理解，为什么这么做？标记+分割号
				//应该是告诉客户端，将要创建的通道
				for i := range tunnelList {
					str += tunnelList[i].Remark + common.CONN_DATA_SEQ
				}

				//这里为什么要写到连接里边
				//database.GetDb().JsonDb.Tasks.Range(func(key, value interface{}) bool {
				//	tun := value.(*models.Tunnel)
				//	//if _, ok := s.runList[v.Id]; ok && v.Client.Id == id {
				//	if _, ok := s.runList.Load(tun.Id); ok && tun.Client.Id == id {
				//		str += tun.Name + common.CONN_DATA_SEQ
				//	}
				//	return true
				//})
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
				if err = database.GetDb().NewClient(client); err != nil {
					fail = true
					c.WriteAddFail()
					break loop
				}
				c.WriteAddOk()
				c.Write([]byte(client.VerifyKey)) //这是为什么？还要向客户端写密钥？
				s.Client.Store(client.Id, NewClient(nil, nil, nil, ""))
			}
		case common.NEW_TASK:
			if tun, err := c.GetTunnelInfo(); err != nil {
				fail = true
				c.WriteAddFail()
				break loop
			} else {
				ports := common.GetPorts(tun.Ports)
				//logs.Info(ports)
				if len(ports) == 0 {
					break loop
				}
				logs.Info(tun.Target.TargetStr)
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
						tunnel.Remark = tun.Remark + "_" + strconv.Itoa(tunnel.ServerPort)
						tunnel.Target = new(models.Target)
						if tun.TargetAddr != "" {
							tunnel.Target.TargetStr = tun.TargetAddr + ":" + strconv.Itoa(targets[i])
						} else {
							tunnel.Target.TargetStr = strconv.Itoa(targets[i])
						}
					}
					//获取新的ID
					//tunnel.Id = database.GetDb().JsonDb.GetTaskId()
					tunnel.Status = true
					tunnel.Flow = new(models.Flow)
					tunnel.NoStore = true
					tunnel.Client = client
					tunnel.Password = tun.Password
					tunnel.LocalPath = tun.LocalPath
					tunnel.StripPre = tun.StripPre
					tunnel.MultiAccount = tun.MultiAccount
					//检查某客户端是否有存在的通道
					if !database.GetDb().HasTunnel(client.Id, tunnel) {
						if err := database.GetDb().NewTask(tunnel); err != nil {
							logs.Notice("Add task error ", err.Error())
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
