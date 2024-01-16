package main

import (
	"context"
	"flag"
	"github.com/astaxie/beego/logs"
	"github.com/kardianos/service"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
	"xnps/client"
	"xnps/database/models"
	"xnps/lib/common"
	"xnps/lib/config"
	"xnps/lib/install"
	"xnps/lib/version"
)

// 获取输入参数
var (
	serverAddr     = flag.String("server", "", "Server addr (ip:port)")
	configPath     = flag.String("config", "", "Configuration file path")
	verifyKey      = flag.String("vkey", "", "Authentication key")
	connType       = flag.String("type", "tcp", "Connection type with the server（kcp|tcp）")
	proxyUrl       = flag.String("proxy", "", "proxy socks5 url(eg:socks5://111:222@127.0.0.1:9007)")
	logLevel       = flag.String("log_level", "7", "log level 0~7")
	registerTime   = flag.Int("time", 2, "register time long /h")
	localPort      = flag.Int("local_port", 2000, "p2p local port")
	password       = flag.String("password", "", "p2p password flag")
	target         = flag.String("target", "", "p2p target")
	localType      = flag.String("local_type", "p2p", "p2p target")
	logPath        = flag.String("log_path", "", "npc log path")
	debug          = flag.Bool("debug", true, "npc debug")
	pprofAddr      = flag.String("pprof", "", "PProf debug addr (ip:port)")
	ver            = flag.Bool("version", false, "show current version")
	disconnectTime = flag.Int("disconnect_timeout", 60, "not receiving check packet times, until timeout will disconnect the client")
)

func main() {
	flag.Parse()
	logs.Reset()
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	if *ver {
		common.PrintVersion()
		return
	}
	if *logPath == "" {
		*logPath = common.GetNpcLogPath()
	}
	if common.IsWindows() {
		*logPath = strings.Replace(*logPath, "\\", "\\\\", -1)
	}
	if *debug {
		logs.SetLogger(logs.AdapterConsole, `{"level":`+*logLevel+`,"color":true}`)
	} else {
		logs.SetLogger(logs.AdapterFile, `{"level":`+*logLevel+`,"filename":"`+*logPath+`","daily":false,"maxlines":100000,"color":true}`)
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
	if err != nil { //输入参数为空时
		slog.Error("npc.go", err, "service function disabled")
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
					slog.Error("npc.go", context.Background())
					log.Println("文件不存在", path)
				}
				client.GetTaskStatus(path)
			}
		//case "register": //不允许注册
		//	flag.CommandLine.Parse(os.Args[2:])
		//	client.RegisterLocalIp(*serverAddr, *verifyKey, *connType, *proxyUrl, *registerTime)
		case "update":
			install.UpdateNpc()
			return
		//case "nat":
		//	c := stun.CreateNewClient()
		//	c.SetServerAddr(*stunAddr)
		//	nat, host, err := c.Discover()
		//	if err != nil || host == nil {
		//		logs.Error("get nat type error", err)
		//		return
		//	}
		//	fmt.Printf("nat type: %s \npublic address: %s\n", nat.String(), host.String())
		//	os.Exit(0)
		case "start", "stop", "restart":
			// support busyBox and sysV, for openWrt
			if service.Platform() == "unix-systemv" {
				logs.Info("unix-systemv service")

				if err = exec.Command("/etc/init.d/"+svcConfig.Name, os.Args[1]).Run(); err != nil {
					slog.Error("npc.go", "service"+os.Args[1]+" err", err)
				}
				return
			}
			err = service.Control(s, os.Args[1])
			if err != nil {
				slog.Error("npc.go", "Valid actions err", os.Args[1], service.ControlAction, err)
			}
			return
		case "install":
			service.Control(s, "stop")
			service.Control(s, "uninstall")
			install.InstallNpc()

			if err = service.Control(s, os.Args[1]); err != nil {
				slog.Error("npc.go", "service"+os.Args[1]+" err", err)
			}
			if service.Platform() == "unix-systemv" {
				logs.Info("unix-systemv service")
				confPath := "/etc/init.d/" + svcConfig.Name
				os.Symlink(confPath, "/etc/rc.d/S90"+svcConfig.Name)
				os.Symlink(confPath, "/etc/rc.d/K02"+svcConfig.Name)
			}
			return
		case "uninstall":
			if err = service.Control(s, os.Args[1]); err != nil {
				slog.Error("npc.go", "service"+os.Args[1]+" err", err)
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
		slog.Error("run npc false", "err", err)
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
			logs.Warning("npc: panic serving %v: %v\n%s", err, string(buf))
		}
	}()
	run()
	select {
	case <-p.exit:
		logs.Warning("stop...")
	}
	return nil
}

func run() {
	//性能分析工具
	common.InitPProfFromArg(*pprofAddr)
	//p2p or secret command
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
		//commonConfig.Client.Cnf = new(models.Config)
		go client.StartLocalServer(localServer, commonConfig)
		return
	}
	env := common.GetEnvMap()
	if *serverAddr == "" {
		*serverAddr, _ = env["NPC_SERVER_ADDR"]
	}
	if *verifyKey == "" {
		*verifyKey, _ = env["NPC_SERVER_VKEY"]
	}
	logs.Info("the version of client is %s, the core version of client is %s", version.VERSION, version.GetCoreVersion())
	if *verifyKey != "" && *serverAddr != "" && *configPath == "" {
		//main
		go func() {
			for {
				client.NewRPClient(*serverAddr, *verifyKey, *connType, *proxyUrl, nil, *disconnectTime).Start()
				logs.Info("Client closed! It will be reconnected in 5 seconds")
				time.Sleep(time.Second * 5)
			}
		}()
	} else {
		if *configPath == "" {
			*configPath = common.GetConfigPath()
		}
		go client.StartFromFile(*configPath)
	}
}
