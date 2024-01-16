package server

import (
	"errors"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
	"xnps/database/Mapper"
	"xnps/database/models"
	"xnps/lib/version"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"xnps/bridge"
	"xnps/lib/common"
	"xnps/server/proxy"
	"xnps/server/tool"
)

var (
	Bridge  *bridge.Bridge
	RunList sync.Map //map[int]interface{}
)

func init() {
	RunList = sync.Map{}
}

// 从数据库初始化通道
// init task from db
func InitFromCsv() {
	//Add a public password
	if vkey := beego.AppConfig.String("public_vkey"); vkey != "" {
		c := Mapper.NewClient(vkey)
		Mapper.GetDb().CreateNewClient(c)
		RunList.Store(c.Id, nil)
		//RunList[c.Id] = nil
	}
	tunList, _ := Mapper.GetDb().GetAllTunnelList(1)
	//Initialize services in server-side files
	for i := range tunList {
		if tunList[i].Valid {
			AddTask(&tunList[i])
		}
	}
	//database.GetDb().JsonDb.Tasks.Range(func(key, value interface{}) bool {
	//	if value.(*models.Tunnel).Status {
	//		AddTask(value.(*models.Tunnel))
	//	}
	//	return true
	//})
}

// DealBridgeTask get bridge command
func DealBridgeTask() {
	for {
		select {
		case t := <-Bridge.OpenTask:
			AddTask(t)
		case t := <-Bridge.CloseTask:
			StopServer(t.Id)
		case id := <-Bridge.CloseClient:
			//DelTunnelAndHostByClientId(id, true)
			err := Mapper.GetDb().DelClient(id)
			if err != nil {
				logs.Warn("del client error, the client don`t exits")
			}

			//if v, ok := database.GetDb().JsonDb.Clients.Load(id); ok {
			//	if v.(*models.Client).Valid {
			//		database.GetDb().DelClient(id)
			//	}
			//}
		case tunnel := <-Bridge.OpenTask:
			StartTask(tunnel.Id)
		case s := <-Bridge.SecretChan:
			logs.Trace("New secret connection, addr", s.Conn.Conn.RemoteAddr())
			if t := Mapper.GetDb().GetTunnelByMd5Password(s.Password); t != nil {
				if t.Status {
					go proxy.NewBaseServer(Bridge, t).DealClient(s.Conn, t.Client, t.Target.TargetStr, nil, common.CONN_TCP, nil, t.Target.LocalProxy, nil)
				} else {
					s.Conn.Close()
					logs.Trace("This key %s cannot be processed,status is close", s.Password)
				}
			} else {
				logs.Trace("This key %s cannot be processed", s.Password)
				s.Conn.Close()
			}
		}
	}
}

// start a new server
func StartNewServer(bridgePort int, cnf *models.Tunnel, bridgeType string, timeout int) {
	Bridge = bridge.NewTunnel(bridgePort, bridgeType, common.GetBoolByStr(beego.AppConfig.String("ip_limit")), RunList, timeout)
	go func() {
		if err := Bridge.StartTunnel(); err != nil {
			logs.Error("start server bridge error", err)
			os.Exit(0)
		}
	}()
	//p2p可以去掉
	//if p, err := beego.AppConfig.Int("p2p_port"); err == nil {
	//	go proxy.NewP2PServer(p).Start()
	//	go proxy.NewP2PServer(p + 1).Start()
	//	go proxy.NewP2PServer(p + 2).Start()
	//}
	go DealBridgeTask()
	go dealClientFlow()
	if svr := NewMode(Bridge, cnf); svr != nil {
		if err := svr.Start(); err != nil {
			logs.Error(err)
		}
		RunList.Store(cnf.Id, svr)
		//RunList[cnf.Id] = svr
	} else {
		logs.Error("Incorrect startup mode %s", cnf.Mode)
	}
}

// 处理终端的流量问题，定时处理
func dealClientFlow() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			dealClientData()
		}
	}
}

// new a server by mode name
func NewMode(Bridge *bridge.Bridge, c *models.Tunnel) proxy.Service {
	var service proxy.Service
	switch c.Mode {
	case "tcp", "file": //这里需要修改随机
		service = proxy.NewTunnelModeServer(proxy.ProcessTunnel, Bridge, c)

	case "tcpTrans":
		service = proxy.NewTunnelModeServer(proxy.HandleTrans, Bridge, c)
	case "udp":
		service = proxy.NewUdpModeServer(Bridge, c)
	case "webServer":
		InitFromCsv()
		t := &models.Tunnel{
			ServerPort: 0,
			Mode:       "httpHostServer",
			Status:     true,
		}
		AddTask(t)
		service = proxy.NewWebServer(Bridge)

	}
	return service
}

// stop server
func StopServer(id int64) error {
	//if v, ok := RunList[id]; ok {
	if v, ok := RunList.Load(id); ok {
		if svr, ok := v.(proxy.Service); ok {
			if err := svr.Close(); err != nil {
				return err
			}
			logs.Info("stop server id %d", id)
		} else {
			logs.Warn("stop server id %d error", id)
		}
		if t, err := Mapper.GetDb().GetTaskById(id); err != nil {
			return err
		} else {
			t.Status = false
			logs.Info("close port %d,remark %s,client id %d,task id %d", t.ServerPort, t.Remark, t.Client.Id, t.Id)
			Mapper.GetDb().UpdateTunnel(t)
		}
		//delete(RunList, id)
		RunList.Delete(id)
		return nil
	}
	return errors.New("task is not running")
}

// add task
func AddTask(t *models.Tunnel) error {
	if t.Mode == "secret" || t.Mode == "p2p" {
		logs.Info("secret task %s start ", t.Remark)
		//RunList[t.Id] = nil
		RunList.Store(t.Id, nil)
		return nil
	}
	if b := tool.TestServerPort(t.ServerPort, t.Mode); !b && t.Mode != "httpHostServer" {
		logs.Error("taskId %d start error port %d open failed", t.Id, t.ServerPort)
		return errors.New("the port open error")
	}
	if minute, err := beego.AppConfig.Int("flow_store_interval"); err == nil && minute > 0 {
		go flowSession(time.Minute * time.Duration(minute))
	}
	if svr := NewMode(Bridge, t); svr != nil {
		logs.Info("tunnel task %s start mode：%s port %d", t.Remark, t.Mode, t.ServerPort)
		//RunList[t.Id] = svr
		RunList.Store(t.Id, svr)
		go func() {
			if err := svr.Start(); err != nil {
				logs.Error("clientId %d taskId %d start error %s", t.Client.Id, t.Id, err)
				//delete(RunList, t.Id)
				RunList.Delete(t.Id)
				return
			}
		}()
	} else {
		return errors.New("the mode is not correct")
	}
	return nil
}

// start task
func StartTask(id int64) error {
	if t, err := Mapper.GetDb().GetTaskById(id); err != nil {
		return err
	} else {
		AddTask(t)
		t.Status = true
		Mapper.GetDb().UpdateTunnel(t)
	}
	return nil
}

// delete task
func DelTask(id int64) error {
	//if _, ok := RunList[id]; ok {
	if _, ok := RunList.Load(id); ok {
		if err := StopServer(id); err != nil {
			return err
		}
	}
	return Mapper.GetDb().DelTunnel(id)
}

// 隧道列表分页，然后返回隧道
// get task list by page num
func GetTunnel(start, length int, modelType string, clientId int64, search string) ([]*models.Tunnel, int) {
	list := make([]*models.Tunnel, 0)
	//var cnt int
	mds, length := Mapper.GetDb().GetTunnelListByClientIdWithPage(start, length, modelType, clientId)
	for c := range mds {
		list = append(list, &mds[c])
	}
	return list, length
	//keys := database.GetMapKeys(database.GetDb().JsonDb.Tasks, false, "", "")
	////这里在遍历所有的隧道的key
	//for _, key := range keys {
	//	if value, ok := database.GetDb().JsonDb.Tasks.Load(key); ok {
	//		v := value.(*models.Tunnel)
	//		if (modelType != "" && v.Mode != modelType || (clientId != 0 && v.Client.Id != clientId)) || (modelType == "" && clientId != v.Client.Id) {
	//			continue
	//		}
	//		if search != "" && !(v.Id == int64(common.GetIntNoErrByStr(search)) || v.ServerPort == common.GetIntNoErrByStr(search) || strings.Contains(v.Password, search) || strings.Contains(v.Name, search) || strings.Contains(v.Client.AccessKey, search)) {
	//			continue
	//		}
	//		cnt++
	//		if _, ok := Bridge.Client.Load(v.Client.Id); ok {
	//			v.Client.Connected = true
	//		} else {
	//			v.Client.Connected = false
	//		}
	//		if start--; start < 0 {
	//			if length--; length >= 0 {
	//				//if _, ok := RunList[v.Id]; ok {
	//				if _, ok := RunList.Load(v.Id); ok {
	//					v.RunStatus = true
	//				} else {
	//					v.RunStatus = false
	//				}
	//				list = append(list, v)
	//			}
	//		}
	//	}
	//}
	//return list, cnt
}

// get client list
func GetClientList(start, length int64, search, sort, order string, clientId int) (list []models.Client, cnt int) {
	list, cnt = Mapper.GetDb().GetAllClientList(start, length, search, sort, order, clientId)
	dealClientData()
	return
}

// TODO:需要重构，添加流量处理相关逻辑
// 处理客户端数据
func dealClientData() {
	//logs.Info("dealClientData.........")
	//这个地方判断为什么？
	Mapper.GetDb().SetClientStatus(true, 0)

	//database.GetDb().JsonDb.Clients.Range(func(key, value interface{}) bool {
	//	v := value.(*models.Client)
	//	if vv, ok := Bridge.Client.Load(v.Id); ok {
	//		v.Connected = true
	//		v.Version = vv.(*bridge.Client).Version
	//	} else {
	//		v.Connected = false
	//	}
	//v.Flow.InletFlow = 0
	//v.Flow.ExportFlow = 0
	//if len(file.GetDb().JsonDb.Hosts) == 0 {
	//
	//}
	//var inflow int64 = 0
	//var outflow int64 = 0
	//file.GetDb().JsonDb.Hosts.Range(func(key, value interface{}) bool {
	//	h := value.(*file.Host)
	//	if h.Client.Id == v.Id {
	//		inflow  += h.Flow.InletFlow
	//		outflow += h.Flow.ExportFlow
	//	}
	//	return true
	//})
	//file.GetDb().JsonDb.Tasks.Range(func(key, value interface{}) bool {
	//	t := value.(*file.Tunnel)
	//	if t.Client.Id == v.Id {
	//		inflow  += t.Flow.InletFlow
	//		outflow += t.Flow.ExportFlow
	//	}
	//	return true
	//})
	//
	//if inflow >0 || outflow >0{
	//	v.Flow.InletFlow = inflow
	//	v.Flow.ExportFlow = outflow
	//}
	//return true
	//})
	return
}

//
//// delete all host and tasks by client id
//func DelTunnelAndHostByClientId(clientId int, justDelNoStore bool) {
//	var ids []int
//	file.GetDb().JsonDb.Tasks.Range(func(key, value interface{}) bool {
//		v := value.(*file.Tunnel)
//		if justDelNoStore && !v.NoStore {
//			return true
//		}
//		if v.Client.Id == clientId {
//			ids = append(ids, v.Id)
//		}
//		return true
//	})
//	for _, id := range ids {
//		DelTunnel(id)
//	}
//	ids = ids[:0]
//	file.GetDb().JsonDb.Hosts.Range(func(key, value interface{}) bool {
//		v := value.(*file.Host)
//		if justDelNoStore && !v.NoStore {
//			return true
//		}
//		if v.Client.Id == clientId {
//			ids = append(ids, v.Id)
//		}
//		return true
//	})
//	for _, id := range ids {
//		file.GetDb().DelHost(id)
//	}
//}

// close the client
func DelClientConnect(clientId int64) {
	Bridge.DelClient(clientId)
}

// 获取后台状态信息
func GetDashboardData() map[string]interface{} {
	data := make(map[string]interface{})
	data["version"] = version.VERSION
	//data["hostCount"] = common.GeSynctMapLen(file.GetDb().JsonDb.Hosts)
	//data["clientCount"] = common.GeSynctMapLen(database.GetDb().JsonDb.Clients)
	data["clientCount"] = Mapper.GetDb().GetAllClientCount(-1, -1)
	if beego.AppConfig.String("public_vkey") != "" { //remove public vkey
		data["clientCount"] = data["clientCount"].(int) - 1
	}
	dealClientData()
	//c := 0
	var in, out int64
	//从数据库获取流量信息
	//获取所有已连接的客户端数量
	//database.GetDb().JsonDb.Clients.Range(func(key, value interface{}) bool {
	//	v := value.(*models.Client)
	//	if v.Connected {
	//		c += 1
	//	}
	//	in += v.Flow.InletFlow
	//	out += v.Flow.ExportFlow
	//	return true
	//})
	//TODO:这里需要检查一下,valid的使用
	data["clientOnlineCount"] = Mapper.GetDb().GetAllClientCount(1, -1)
	data["inletFlowCount"] = int(in)
	data["exportFlowCount"] = int(out)
	//分别获得各类隧道的数量
	var tcp = Mapper.GetDb().GetClientCountByMode("tcp")
	var udp = Mapper.GetDb().GetClientCountByMode("udp")
	var secret = Mapper.GetDb().GetClientCountByMode("secret")
	var socks5 = Mapper.GetDb().GetClientCountByMode("socks5")
	var p2p = Mapper.GetDb().GetClientCountByMode("p2p")
	var http = Mapper.GetDb().GetClientCountByMode("httpProxy")
	//var udp, secret, socks5, p2p, http int
	//
	//database.GetDb().JsonDb.Tasks.Range(func(key, value interface{}) bool {
	//	switch value.(*models.Tunnel).Mode {
	//	case "tcp":
	//		tcp += 1
	//	case "socks5":
	//		socks5 += 1
	//	case "httpProxy":
	//		http += 1
	//	case "udp":
	//		udp += 1
	//	case "p2p":
	//		p2p += 1
	//	case "secret":
	//		secret += 1
	//	}
	//	return true
	//})

	data["tcpC"] = tcp
	data["udpCount"] = udp
	data["socks5Count"] = socks5
	data["httpProxyCount"] = http
	data["secretCount"] = secret
	data["p2pCount"] = p2p
	data["bridgeType"] = beego.AppConfig.String("bridge_type")
	data["httpProxyPort"] = beego.AppConfig.String("http_proxy_port")
	data["httpsProxyPort"] = beego.AppConfig.String("https_proxy_port")
	data["ipLimit"] = beego.AppConfig.String("ip_limit")
	data["flowStoreInterval"] = beego.AppConfig.String("flow_store_interval")
	data["serverIp"] = beego.AppConfig.String("p2p_ip")
	data["p2pPort"] = beego.AppConfig.String("p2p_port")
	data["logLevel"] = beego.AppConfig.String("log_level")
	//tcpCount := 0
	//
	//database.GetDb().JsonDb.Clients.Range(func(key, value interface{}) bool {
	//	tcpCount += int(value.(*models.Client).NowConn)
	//	return true
	//})
	count := Mapper.GetDb().GetAllTunnelCountByStatus(true, common.MODE_TCP)
	data["tcpCount"] = count
	cpuPercet, _ := cpu.Percent(0, true)
	var cpuAll float64
	for _, v := range cpuPercet {
		cpuAll += v
	}
	loads, _ := load.Avg()
	data["load"] = loads.String()
	data["cpu"] = math.Round(cpuAll / float64(len(cpuPercet)))
	swap, _ := mem.SwapMemory()
	data["swap_mem"] = math.Round(swap.UsedPercent)
	vir, _ := mem.VirtualMemory()
	data["virtual_mem"] = math.Round(vir.UsedPercent)
	conn, _ := net.ProtoCounters(nil)
	io1, _ := net.IOCounters(false)
	time.Sleep(time.Millisecond * 500)
	io2, _ := net.IOCounters(false)
	if len(io2) > 0 && len(io1) > 0 {
		data["io_send"] = (io2[0].BytesSent - io1[0].BytesSent) * 2
		data["io_recv"] = (io2[0].BytesRecv - io1[0].BytesRecv) * 2
	}
	for _, v := range conn {
		data[v.Protocol] = v.Stats["CurrEstab"]
	}
	//chart
	var fg int
	if len(tool.ServerStatus) >= 10 {
		fg = len(tool.ServerStatus) / 10
		for i := 0; i <= 9; i++ {
			data["sys"+strconv.Itoa(i+1)] = tool.ServerStatus[i*fg]
		}
	}
	return data
}

// 实例化流量数据到文件
func flowSession(m time.Duration) {
	//时间触发器，可以学习
	ticker := time.NewTicker(m)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			//file.GetDb().JsonDb.StoreHostToJsonFile()
			//database.GetDb().JsonDb.StoreTasksToJsonFile()
			//database.GetDb().JsonDb.StoreClientsToJsonFile()
		}
	}
}
