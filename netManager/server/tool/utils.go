package tool

import (
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"xnps/lib/common"
)

var (
	ports        []uint16
	ServerStatus []map[string]interface{}
)

// 监控服务端状态
func StartSystemInfo() {
	if b, err := beego.AppConfig.Bool("system_info_display"); err == nil && b {
		ServerStatus = make([]map[string]interface{}, 0, 1500)
		go getSeverStatus()
	}
}

func InitAllowPort() {
	p := beego.AppConfig.String("allow_ports")
	ports = common.GetPorts(p)
}

func TestServerPort(port uint16, m string) (b bool) {
	if m == "p2p" || m == "secret" {
		return true
	}
	if port > 65535 || port < 0 {
		return false
	}
	if len(ports) != 0 {
		if !common.InIntArr(ports, port) {
			return false
		}
	}
	if m == "udp" {
		b = common.TestUdpPort(port)
	} else {
		b = common.TestTcpPort(port)
	}
	return
}

func StrToInt64(s string, def ...int64) int64 {
	out, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return def[0]
	}
	return out
}

// TODO:设置端口范围 端口不能这样使用随机数生成，应该在某一范围内随机顺序生成
func GenerateServerPort(m string) uint16 {
	for {
		//生成随机数 1024 - 65535
		i := rand.Intn(len(ports))
		serverPort := ports[i]

		if TestServerPort(serverPort, m) {
			return serverPort
		}
	}
}

func getSeverStatus() {
	for {
		if len(ServerStatus) < 10 {
			time.Sleep(time.Second)
		} else {
			time.Sleep(time.Minute)
		}
		cpuPercet, _ := cpu.Percent(0, true)
		var cpuAll float64
		for _, v := range cpuPercet {
			cpuAll += v
		}
		m := make(map[string]interface{})
		loads, _ := load.Avg()
		m["load1"] = loads.Load1
		m["load5"] = loads.Load5
		m["load15"] = loads.Load15
		m["cpu"] = math.Round(cpuAll / float64(len(cpuPercet)))
		swap, _ := mem.SwapMemory()
		m["swap_mem"] = math.Round(swap.UsedPercent)
		vir, _ := mem.VirtualMemory()
		m["virtual_mem"] = math.Round(vir.UsedPercent)
		conn, _ := net.ProtoCounters(nil)
		io1, _ := net.IOCounters(false)
		time.Sleep(time.Millisecond * 500)
		io2, _ := net.IOCounters(false)
		if len(io2) > 0 && len(io1) > 0 {
			m["io_send"] = (io2[0].BytesSent - io1[0].BytesSent) * 2
			m["io_recv"] = (io2[0].BytesRecv - io1[0].BytesRecv) * 2
		}
		t := time.Now()
		m["time"] = strconv.Itoa(t.Hour()) + ":" + strconv.Itoa(t.Minute()) + ":" + strconv.Itoa(t.Second())

		for _, v := range conn {
			m[v.Protocol] = v.Stats["CurrEstab"]
		}
		if len(ServerStatus) >= 1440 {
			ServerStatus = ServerStatus[1:]
		}
		ServerStatus = append(ServerStatus, m)
	}
}
