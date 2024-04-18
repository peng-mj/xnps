package config

import (
	"os"
	"xnps/pkg/sysTool"

	"golang.org/x/exp/slog"
)

var (
	sysConf SysConfig
)

const (
	SoftWareVersion = "XNPS_2023-0.0.1"
)

type SysConfig struct {
	driver *Driver `json:"-"`
}

func SysConfigInit() {
	var err error
	if sysTool.FileExisted("/etc/xnps/xnps.ini") {
		sysConf.driver, err = NewSysConfig("/etc/xnps/xnps.ini")
	} else if sysTool.FileExisted("./conf/xnps.ini") {
		sysConf.driver, err = NewSysConfig("./conf/xnps.ini")
	} else {
		slog.Info(``)
		sysTool.CreateFolder("./conf")
		sysTool.CreateAndWriteFile("./conf/xnps.ini", InitFileContent)
		sysConf.driver, err = NewSysConfig("./conf/xnps.ini")
	}
	if err != nil {
		slog.Info("打开配置文件失败,请检查")
		os.Exit(-1)
	}

}

func SysConf() *SysConfig {
	return &sysConf
}
