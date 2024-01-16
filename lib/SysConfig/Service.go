package SysConfig

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"xnps/lib/SysTool"
	"xnps/lib/crypt"
)

var sysConf *SysConfig

func (s *SysConfig) CheckUserName(username string) bool {
	if len(username) > 3 {
		return username == s.WebUserName
	}
	return false
}

func SysConfigInit(path string) {
	var err error

	if SysTool.FileExisted(path) {
		path = "/etc/sunrun/sunrun.ini"
		sysConf.driver, err = NewSysConfig(path)
	} else if SysTool.FileExisted("./conf/sunrun.ini") {
		path = "./conf/sunrun.ini"
		sysConf.driver, err = NewSysConfig("./conf/sunrun.ini")
	} else {
		slog.Error(`
	配置文件读取错误
请检查路径：/etc/sunrun/sunrun.confDriver 或者 ./confDriver/sunrun.conf下文件是否存在
将在当前文件夹下创建默认配置文件`)
		if runtime.GOOS == "linux" {
			path = "/etc/sunrun/"
		} else {
			path = "./conf/"
		}
		SysTool.CreateFolder(path)
		path = filepath.Join(path, "sunrun.ini")
		SysTool.CreateAndWriteFile(path, InitFileContentNPS)
		sysConf.driver, err = NewSysConfig(path)
	}
	if err != nil {
		slog.Info("打开配置文件失败,请检查")
		os.Exit(-1)
	}

	sysConf.WebHost = sysConf.driver.GetValue("sys", "server_host")
	if len(sysConf.WebHost) < 10 {
		sysConf.WebHost = "127.0.0.1:8888"
		sysConf.driver.SetValue("sys", "server_host", sysConf.WebHost)
	}
	sysConf.AuthCryptKey = sysConf.driver.GetValue("sys", "client_uuid")
	if len(sysConf.AuthCryptKey) < 6 {
		sysConf.AuthCryptKey = crypt.GetRandomString(30)
		sysConf.driver.SetValue("sys", "client_uuid", sysConf.AuthCryptKey)
	}

	sysConf.WebPort = sysConf.driver.GetInt64("web", "port", 9870)
	if sysConf.WebPort > 65535 && sysConf.WebPort < 80 {
		sysConf.WebPort = 9870
	}

	//slog.Info("设备状态信息", status)
	slog.Info("设备配置信息", sysConf)
}
