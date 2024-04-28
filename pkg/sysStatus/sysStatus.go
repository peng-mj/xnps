package sysStatus

import (
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"net"
	"runtime"
	"time"
)

type SysStatus struct {
	SysTime   int64  `json:"sysTime"`
	TotalMem  uint64 `json:"totalMem"`
	UsedMen   uint64 `json:"UsedMen"`
	TotalDisk uint64 `json:"TotalDisk"`
	UsedDisk  uint64 `json:"usedDisk"`
	BootTime  uint64 `json:"BootTime"`
	UpTime    uint64 `json:"upTime"`
	Ip        string `json:"ip"`
}

func (d *SysStatus) GetSysStatus() (status SysStatus) {
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
	// fmt.Println(localAddr.String())
	return RemoteAddr.IP.String(), LanAddr.IP.String()
}
