package main

import (
	"flag"
	"github.com/astaxie/beego"
	"github.com/kardianos/service"
	"golang.org/x/exp/slog"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"tunpx/pkg/common"
	"tunpx/pkg/daemon"
	"tunpx/pkg/install"
	"tunpx/pkg/sysTool"
)

var (
	level string
	ver   = flag.Bool("version", false, "show current version")
)

func main() {

	flag.Parse()
	// init log
	if *ver {
		common.PrintVersion()
		return
	}

	// *confPath why get null value ?
	for _, v := range os.Args[1:] {
		switch v {
		case "install", "start", "stop", "uninstall", "restart":
			continue
		}
		if strings.Contains(v, "-conf_path=") {
			common.ConfPath = strings.Replace(v, "-conf_path=", "", -1)
		}
	}

	if err := beego.LoadAppConfig("ini", filepath.Join(common.GetRunPath(), "conf", "tunpxs.conf")); err != nil {
		slog.Error("加载配置文件失败", "load config file error", err)
	}

	common.InitPProfFromFile()
	if level = beego.AppConfig.String("log_level"); level == "" {
		level = "7"
	}
	logPath := "./log/tunpxs.log"
	if !sysTool.DirExisted("./log") {
		sysTool.CreateFolder("./log")
	}

	r := &lumberjack.Logger{
		Filename:   logPath,
		LocalTime:  true,
		MaxSize:    20,
		MaxAge:     7,
		MaxBackups: 7,
		Compress:   true,
	}
	logger := slog.New(slog.NewJSONHandler(r, nil))
	slog.SetDefault(logger)
	// init service
	options := make(service.KeyValue)
	svcConfig := &service.Config{
		Name:        "xNps",
		DisplayName: "xnps内网穿透代理服务器",
		Description: "一款轻量级、功能强大的内网穿透代理服务器。支持tcp、udp流量转发，支持内网http代理、内网socks5代理，同时支持snappy压缩、站点保护、加密传输、多路复用、header修改等。支持web图形化管理，集成多用户模式。",
		Option:      options,
	}

	for _, v := range os.Args[1:] {
		switch v {
		case "install", "start", "stop", "uninstall", "restart":
			continue
		}
		svcConfig.Arguments = append(svcConfig.Arguments, v)
	}

	svcConfig.Arguments = append(svcConfig.Arguments, "service")

	if !common.IsWindows() {
		svcConfig.Dependencies = []string{
			"Requires=network.target",
			"After=network-online.target syslog.target"}
		svcConfig.Option["SystemdScript"] = install.SystemdScript
		svcConfig.Option["SysvScript"] = install.SysvScript
	}
	prg := &nps{}
	prg.exit = make(chan struct{})
	s, err := service.New(prg, svcConfig)
	if err != nil {
		slog.Error("service function disabled", "error", err)
		run()
		// run without service
		wg := sync.WaitGroup{}
		wg.Add(1)
		wg.Wait()
		return
	}

	if len(os.Args) > 1 && os.Args[1] != "service" {
		switch os.Args[1] {
		case "reload":
			daemon.InitDaemon("tunpxs", common.GetRunPath(), common.GetTmpPath())
			return
		case "install":
			// uninstall before
			_ = service.Control(s, "stop")
			_ = service.Control(s, "uninstall")

			binPath := install.InstallNps()
			svcConfig.Executable = binPath
			s, err = service.New(prg, svcConfig)
			if err != nil {
				slog.Error("create service error", "error", err)
				return
			}
			err = service.Control(s, os.Args[1])
			if err != nil {
				slog.Error("Valid actions: %q\n%s", service.ControlAction, err.Error())
			}
			if service.Platform() == "unix-systemv" {
				slog.Info("unix-systemv service")
				confPath := "/etc/init.d/" + svcConfig.Name
				os.Symlink(confPath, "/etc/rc.d/S90"+svcConfig.Name)
				os.Symlink(confPath, "/etc/rc.d/K02"+svcConfig.Name)
			}
			return
		case "start", "restart", "stop":
			if service.Platform() == "unix-systemv" {
				slog.Info("unix-systemv service")
				err = exec.Command("/etc/init.d/"+svcConfig.Name, os.Args[1]).Run()
				if err != nil {
					slog.Error("运行系统命令出错", "cmd", os.Args[1], "error", err)
				}
				return
			}
			err := service.Control(s, os.Args[1])
			if err != nil {
				slog.Error("Valid actions: %q\n%s", service.ControlAction, err.Error())
			}
			return
		case "uninstall":
			err := service.Control(s, os.Args[1])
			if err != nil {
				slog.Error("Valid actions: %q\n%s", service.ControlAction, err.Error())
			}
			if service.Platform() == "unix-systemv" {
				slog.Info("unix-systemv service")
				os.Remove("/etc/rc.d/S90" + svcConfig.Name)
				os.Remove("/etc/rc.d/K02" + svcConfig.Name)
			}
			return
		case "update":
			install.UpdateNps()
			return
			// default:
			//	slog.Error("command is not support")
			//	return
		}
	}

	_ = s.Run()
}

type nps struct {
	exit chan struct{}
}

func (p *nps) Start(s service.Service) error {
	_, _ = s.Status()
	go p.run()
	return nil
}
func (p *nps) Stop(s service.Service) error {
	_, _ = s.Status()
	close(p.exit)
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func (p *nps) run() error {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			slog.Warn("tunpxs: panic serving", err, string(buf))
		}
	}()
	run()
	select {
	case <-p.exit:
		slog.Warn("stop...")
	}
	return nil
}

func run() {
	// task := &models.Tunnel{
	//	Mode: "webServer",
	// }
	// bridgePort, err := beego.AppConfig.Int("bridge_port")
	// if err != nil {
	//	slog.Error("Getting bridge_port error", err)
	//	os.Exit(0)
	// }
	//
	// slog.Info("the config path is:" + common.GetRunPath())
	// slog.Info("the version of server is %s ,allow client core version to be %s", version.VERSION, version.GetCoreVersion())
	// //初始化
	// connection.InitConnectionService()
	// //crypt.InitTls(filepath.Join(common.GetRunPath(), "conf", "server.pem"), filepath.Join(common.GetRunPath(), "conf", "server.key"))
	// //初始化密钥
	// crypt.InitTls()
	// //初始化允许的端口
	// tool.InitAllowPort()
	// //持续获取系统信息
	// tool.StartSystemInfo()
	// //
	// timeout, err := beego.AppConfig.Int("disconnect_timeout")
	// if err != nil {
	//	timeout = 60
	// }
	// go server.StartNewServer(bridgePort, task, beego.AppConfig.String("bridge_type"), timeout)
}
