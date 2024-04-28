package main

import (
	"flag"
	"github.com/kardianos/service"
	"golang.org/x/exp/slog"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
	"tunpx/pkg/common"
	"tunpx/pkg/config"
	_ "tunpx/pkg/crypt"
	"tunpx/pkg/install"
	"tunpx/pkg/models"
	"tunpx/pkg/sysTool"
	"tunpx/pkg/version"
)

// 获取输入参数
var (
	serverAddr     = flag.String("server", "", "Server addr (ip:port)")
	configPath     = flag.String("config", "", "Configuration file path")
	verifyKey      = flag.String("vkey", "", "Authentication key")
	connType       = flag.String("type", "tcp", "Connection type with the server（kcp|tcp）")
	proxyUrl       = flag.String("proxy", "", "proxy socks5 url(eg:socks5://111:222@127.0.0.1:9007)")
	localPort      = flag.Int("local_port", 2000, "p2p local port")
	password       = flag.String("password", "", "p2p password flag")
	target         = flag.String("target", "", "p2p target")
	localType      = flag.String("local_type", "p2p", "p2p target")
	logPath        = flag.String("log_path", "", "tunpxc log path")
	pprofAddr      = flag.String("pprof", "", "PProf debug addr (ip:port)")
	ver            = flag.Bool("version", false, "show current version")
	disconnectTime = flag.Int("disconnect_timeout", 60, "not receiving check packet times, until timeout will disconnect the client")
)

func main() {
	if *logPath == "" {
		*logPath = "./log/tunpxs.log"
		if !sysTool.DirExisted("./log") {
			sysTool.CreateFolder("./log")
		}
	}

	r := &lumberjack.Logger{
		Filename:   *logPath,
		LocalTime:  true,
		MaxSize:    20,
		MaxAge:     7,
		MaxBackups: 7,
		Compress:   true,
	}
	logger := slog.New(slog.NewJSONHandler(r, nil))
	slog.SetDefault(logger)

	flag.Parse()

	if *ver {
		common.PrintVersion()
		return
	}

	// init service
	options := make(service.KeyValue)
	svcConfig := &service.Config{
		Name:        "Npc",
		DisplayName: "nps内网穿透客户端",
		Description: "一款轻量级、功能强大的内网穿透代理服务器。",
		Option:      options,
	}
	if !common.IsWindows() {
		svcConfig.Dependencies = []string{
			"Requires=network.target",
			"After=network-online.target syslog.target"}
		svcConfig.Option["SystemdScript"] = install.SystemdScript
		svcConfig.Option["SysvScript"] = install.SysvScript
	}
	for _, v := range os.Args[1:] {
		switch v {
		case "install", "start", "stop", "uninstall", "restart":
			continue
		}
		if !strings.Contains(v, "-service=") && !strings.Contains(v, "-debug=") {
			svcConfig.Arguments = append(svcConfig.Arguments, v)
		}
	}
	svcConfig.Arguments = append(svcConfig.Arguments, "-debug=false")
	prg := &npc{
		exit: make(chan struct{}),
	}
	s, err := service.New(prg, svcConfig)
	if err != nil { // 输入参数为空时
		slog.Error("xnpc.go", err, "service function disabled")
		run()
		// run without service
		wg := sync.WaitGroup{}
		wg.Add(1)
		wg.Wait()
		return
	}
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "status":
			if len(os.Args) > 2 {
				path := strings.Replace(os.Args[2], "-config=", "", -1)
				if len(path) < 3 {
					slog.Error("配置文件不存在", "path", path)
				}
				client.GetTaskStatus(path)
			}
		case "update":
			install.UpdateNpc()
			return
		case "start", "stop", "restart":
			// support busyBox and sysV, for openWrt
			if service.Platform() == "unix-systemv" {
				slog.Info("unix-systemv service")

				if err = exec.Command("/etc/init.d/"+svcConfig.Name, os.Args[1]).Run(); err != nil {
					slog.Error("tunpxc.go", "service"+os.Args[1]+" err", err)
				}
				return
			}
			err = service.Control(s, os.Args[1])
			if err != nil {
				slog.Error("tunpxc.go", "Valid actions err", os.Args[1], service.ControlAction, err)
			}
			return
		case "install":
			service.Control(s, "stop")
			service.Control(s, "uninstall")
			install.InstallNpc()

			if err = service.Control(s, os.Args[1]); err != nil {
				slog.Error("tunpxc.go", "service"+os.Args[1]+" err", err)
			}
			if service.Platform() == "unix-systemv" {
				slog.Info("unix-systemv service")
				confPath := "/etc/init.d/" + svcConfig.Name
				os.Symlink(confPath, "/etc/rc.d/S90"+svcConfig.Name)
				os.Symlink(confPath, "/etc/rc.d/K02"+svcConfig.Name)
			}
			return
		case "uninstall":
			if err = service.Control(s, os.Args[1]); err != nil {
				slog.Error("tunpxc.go", "service"+os.Args[1]+" err", err)
			}
			if service.Platform() == "unix-systemv" {
				slog.Info("unix-systemv service")
				os.Remove("/etc/rc.d/S90" + svcConfig.Name)
				os.Remove("/etc/rc.d/K02" + svcConfig.Name)
			}
			return
		}
	}
	if err = s.Run(); err != nil {
		slog.Error("run tunpxc false", "err", err)
	}

}

type npc struct {
	exit chan struct{}
}

func (p *npc) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *npc) Stop(s service.Service) error {
	close(p.exit)
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func (p *npc) run() error {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			slog.Warn("tunpxc: panic serving ", err, string(buf))
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
	// 性能分析工具
	common.InitPProfFromArg(*pprofAddr)
	// p2p or secret command
	if *password != "" {
		commonConfig := new(config.CommonConfig)
		commonConfig.Server = *serverAddr
		commonConfig.VKey = *verifyKey
		commonConfig.Tp = *connType
		localServer := new(config.LocalServer)
		localServer.Type = *localType
		localServer.Password = *password
		localServer.Target = *target
		localServer.Port = *localPort
		commonConfig.Client = new(models.Client)
		// commonConfig.Client.Cnf = new(models.Config)
		go client2.StartLocalServer(localServer, commonConfig)
		return
	}
	env := common.GetEnvMap()
	if *serverAddr == "" {
		*serverAddr, _ = env["NPC_SERVER_ADDR"]
	}
	if *verifyKey == "" {
		*verifyKey, _ = env["NPC_SERVER_VKEY"]
	}
	slog.Info("version info", "client version", version.VERSION, "core version", version.GetCoreVersion())
	if *verifyKey != "" && *serverAddr != "" && *configPath == "" {
		// main
		go func() {
			for {
				client.NewRPClient(*serverAddr, *verifyKey, *connType, *proxyUrl, nil, *disconnectTime).Start()
				slog.Info("Client closed! It will be reconnected in 5 seconds")
				time.Sleep(time.Second * 5)
			}
		}()
	} else {
		if *configPath == "" {
			*configPath = common.GetConfigPath()
		}
		go client2.StartFromFile(*configPath)
	}
}
