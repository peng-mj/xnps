package config

import (
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"log"
	"log/slog"
	"net"
	"os"
	"runtime"
	"time"
	"xnps/lib/SysTool"
)

var (
	sysConf SysConfig
)

const (
	SoftWareVersion = "SUNRUN_CLIENT_KS_2023-0.0.1"
)

type SysStatus struct {
	SysTime    int64  `json:"sysTime"`
	CamaraTime int64  `json:"camaraTime"`
	TotalMem   uint64 `json:"totalMem"`
	UsedMen    uint64 `json:"UsedMen"`
	TotalDisk  uint64 `json:"TotalDisk"`
	UsedDisk   uint64 `json:"usedDisk"`
	BootTime   uint64 `json:"BootTime"`
	UpTime     uint64 `json:"upTime"`
	Ip         string `json:"ip"`
}

type SysConfig struct {
	WebPort    int64   `json:"webPort"`    //设备端后台运行的端口
	WebUser    string  `json:"webUser"`    //设备端后台网页的用户
	WebPasswd  string  `json:"webPasswd"`  //设备端后台网页的密码
	SysVersion string  `json:"sysVersion"` //软件版本
	driver     *Driver `json:"-"`
}

func SysConfigInit() {
	var err error
	if SysTool.FileExisted("/etc/xnps/xnps.ini") {
		sysConf.driver, err = NewSysConfig("/etc/xnps/xnps.ini")
	} else if SysTool.FileExisted("./conf/xnps.ini") {
		sysConf.driver, err = NewSysConfig("./conf/xnps.ini")
	} else {
		slog.Info(`
	配置文件读取错误
请检查路径：/etc/xnps/xnps.ini 或者 ./conf/xnps.ini下文件是否存在
将在当前文件夹下创建默认配置文件`)
		SysTool.CreateFolder("./conf")
		SysTool.CreateAndWriteFile("./conf/xnps.ini", InitFileContent)
		sysConf.driver, err = NewSysConfig("./conf/xnps.ini")
	}
	if err != nil {
		slog.Info("打开配置文件失败,请检查")
		os.Exit(-1)
	}
	sysConf.SysVersion = SoftWareVersion

	sysConf.WebPort = sysConf.driver.GetInt64("web", "port", 8888)
	if sysConf.WebPort > 65535 {
		sysConf.WebPort = 8888
	}

}

func SysConf() *SysConfig {
	return &sysConf
}

func (d *SysConfig) GetSysStatus() (status SysStatus) {
	var bootTime uint64
	memInfo, _ := mem.VirtualMemory()
	if bootTime < 1700000000 {
		bootTime, _ = host.BootTime()
		if bootTime < 1700000000 {
			bootTime = uint64(time.Now().Unix())
		}
	}
	status.BootTime = bootTime
	status.UpTime, _ = host.Uptime()
	status.SysTime = time.Now().Unix()
	status.TotalMem = memInfo.Total / 1000000
	status.UsedMen = memInfo.Used / 1000000
	status.TotalDisk, status.UsedDisk = getDiskSpace()
	_, status.Ip = GetIPInfo()
	return
}
func getDiskSpace() (totalSpace, usedSpace uint64) {
	root := ""
	if runtime.GOOS == "linux" {
		root = "/"
	} else {
		root = "C:"
	}
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Println("Failed to get partitions:", err)
		return
	}
	for _, partition := range partitions {
		if partition.Mountpoint == root {
			usage, err := disk.Usage(partition.Mountpoint)
			if err != nil {
				log.Println("Failed to get partition usage:", err)
				continue
			}
			totalSpace = usage.Total / (1024 * 1024)               // MB
			usedSpace = (usage.Total - usage.Free) / (1024 * 1024) // MB
			return
		}
	}
	return
}
func GetIPInfo() (remote, lan string) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", "127.0.0.1"
	}
	defer conn.Close()
	RemoteAddr := conn.RemoteAddr().(*net.UDPAddr)
	LanAddr := conn.LocalAddr().(*net.UDPAddr)
	//fmt.Println(localAddr.String())
	return RemoteAddr.IP.String(), LanAddr.IP.String()
}
