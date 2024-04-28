package common

import (
	"golang.org/x/exp/slog"
	"os"
	"path/filepath"
	"runtime"
	"tunpx/pkg/sysTool"
)

var ConfPath string

// Get the currently selected configuration file directory
// For non-Windows systems, select the /etc/tunpxs as config directory if exist, or select ./
// windows system, select the C:\Program Files\tunpxs as config directory if exist, or select ./
func GetRunPath() string {
	var path string
	if len(os.Args) == 1 {
		if !IsWindows() {
			dir, _ := filepath.Abs(filepath.Dir(os.Args[0])) // 返回
			return dir + "/"
		} else {
			return "./"
		}
	} else {
		if path = GetInstallPath(); !FileExists(path) {
			return GetAppPath()
		}
	}
	return path
}

// Different systems get different installation paths
func GetInstallPath() string {
	var path string

	if ConfPath != "" {
		return ConfPath
	}

	if IsWindows() {
		path = `C:\Program Files\nps`
	} else {
		path = "/etc/tunpx"
	}

	return path
}

// Get the absolute path to the running directory
func GetAppPath() string {
	if path, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		return path
	}
	return os.Args[0]
}

// Determine whether the current system is a Windows system?
func IsWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}

// interface log file path
func GetLogPath() string {
	var path string
	if IsWindows() {
		path = filepath.Join(GetAppPath(), "tunpxs.log")
	} else {
		path = "/var/log/tunpx.log"
	}
	return path
}

// interface tunpxc log file path
func GetNpcLogPath() string {
	var path string
	if IsWindows() {
		path = filepath.Join(GetAppPath(), "tunpx.log")
	} else {
		path = "/var/log/tunpx.log"
	}
	return path
}

// interface pid file path
func GetTmpPath() string {
	var path string
	if IsWindows() {
		path = GetAppPath()
	} else {
		path = "/tmp"
	}
	return path
}

// config file path
func GetConfigPath() (path string) {
	if IsWindows() {
		path = "conf/tunpx.conf"
		if sysTool.FileExisted(path) {
			return
		} else if path = filepath.Join(GetAppPath(), "conf/tunpx.conf"); sysTool.FileExisted(path) {
			return
		} else {
			slog.Error("加载配置文件失败，请检查", "配置文件状态", "不存在")
			os.Exit(0)
		}
	} else {
		path = "conf/tunpx.conf"
		if sysTool.FileExisted(path) {
			return
		} else if path = "/etc/tunpx/conf/tunpx.conf"; sysTool.FileExisted(path) {
			return
		} else if path = "~/.tunpx/tunpx.conf"; sysTool.FileExisted(path) {
			return
		} else {

		}
		_, err := os.Lstat(path)
		if err != nil {
			path = "/etc/tunpx/conf/tunpx.conf"
		}
	}
	return
}
