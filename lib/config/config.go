package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"xnps/database/models"
	"xnps/lib/common"
)

type CommonConfig struct {
	Server           string
	VKey             string
	Tp               string //bridgeType kcp or tcp
	AutoReconnection bool
	ProxyUrl         string
	Client           *models.Client
	DisconnectTime   int
}

type LocalServer struct {
	Type     string
	Port     int
	Ip       string
	Password string
	Target   string
}

type Config struct {
	content      string
	title        []string
	CommonConfig *CommonConfig
	//Hosts        []*file.Host
	Tasks []*models.Tunnel
	//Healths     []*file.Health
	LocalServer []*LocalServer
}

func NewConfig(path string) (c *Config, err error) {
	c = new(Config)
	var b []byte
	if b, err = common.ReadAllFromFile(path); err != nil {
		return
	} else {
		if c.content, err = common.ParseStr(string(b)); err != nil {
			return nil, err
		}
		if c.title, err = getAllTitle(c.content); err != nil {
			return
		}
		var nowIndex int
		var nextIndex int
		var nowContent string
		for i := 0; i < len(c.title); i++ {
			nowIndex = strings.Index(c.content, c.title[i]) + len(c.title[i])
			if i < len(c.title)-1 {
				nextIndex = strings.Index(c.content, c.title[i+1])
			} else {
				nextIndex = len(c.content)
			}
			nowContent = c.content[nowIndex:nextIndex]

			if strings.Index(getTitleContent(c.title[i]), "secret") == 0 && !strings.Contains(nowContent, "mode") {
				local := delLocalService(nowContent)
				local.Type = "secret"
				c.LocalServer = append(c.LocalServer, local)
				continue
			}
			//except mode
			if strings.Index(getTitleContent(c.title[i]), "p2p") == 0 && !strings.Contains(nowContent, "mode") {
				local := delLocalService(nowContent)
				local.Type = "p2p"
				c.LocalServer = append(c.LocalServer, local)
				continue
			}
			//health set
			//if strings.Index(getTitleContent(c.title[i]), "health") == 0 {
			//	c.Healths = append(c.Healths, dealHealth(nowContent))
			//	continue
			//}
			switch c.title[i] {
			case "[common]":
				c.CommonConfig = dealCommon(nowContent)
			default:
				if strings.Index(nowContent, "host") > -1 {
					//h := dealHost(nowContent)
					//h.Name = getTitleContent(c.title[i])
					//c.Hosts = append(c.Hosts, h)
				} else {
					t := dealTunnel(nowContent)
					t.Remark = getTitleContent(c.title[i])
					c.Tasks = append(c.Tasks, t)
				}
			}
		}
	}
	return
}

func getTitleContent(s string) string {
	re, _ := regexp.Compile(`[\[\]]`)
	return re.ReplaceAllString(s, "")
}

func dealCommon(s string) *CommonConfig {
	c := &CommonConfig{}
	//c.Client = Mapper.NewClient("")
	c.Client = new(models.Client)
	//c.Client.Cnf = new(models.Config)
	for _, v := range splitStr(s) {
		item := strings.Split(v, "=")
		if len(item) == 0 {
			continue
		} else if len(item) == 1 {
			item = append(item, "")
		}
		switch item[0] {
		case "server_addr":
			c.Server = item[1]
		case "vkey":
			c.VKey = item[1]
		case "conn_type":
			c.Tp = item[1]
		case "auto_reconnection":
			c.AutoReconnection = common.GetBoolByStr(item[1])
		case "compress":
			c.Client.Compress = common.GetBoolByStr(item[1])
		case "crypt":
			c.Client.Crypt = common.GetBoolByStr(item[1])
		case "proxy_url":
			c.ProxyUrl = item[1]
		case "rate_limit":
			c.Client.RateLimit = common.GetIntNoErrByStr(item[1])
		case "max_conn":
			c.Client.MaxConn = common.GetIntNoErrByStr(item[1])
		case "remark":
			c.Client.Name = item[1]
		case "pprof_addr":
			common.InitPProfFromArg(item[1])
		case "disconnect_timeout":
			c.DisconnectTime = common.GetIntNoErrByStr(item[1])
		}
	}
	return c
}

func dealTunnel(s string) *models.Tunnel {
	t := &models.Tunnel{}
	t.Target = new(models.Target)
	for _, v := range splitStr(s) {
		item := strings.Split(v, "=")
		if len(item) == 0 {
			continue
		} else if len(item) == 1 {
			item = append(item, "")
		}
		switch strings.TrimSpace(item[0]) {
		case "server_port":
			t.Ports = item[1]
		case "server_ip":
			t.ServerIp = item[1]
		case "mode":
			t.Mode = item[1]
		case "target_addr":
			t.Target.TargetStr = strings.Replace(item[1], ",", "\n", -1)
		case "target_port":
			t.Target.TargetStr = item[1]
		case "target_ip":
			t.TargetAddr = item[1]
		case "password":
			t.Password = item[1]
		case "local_path":
			t.LocalPath = item[1]
		case "strip_pre":
			t.StripPre = item[1]

		}
	}
	return t

}

func delLocalService(s string) *LocalServer {
	l := new(LocalServer)
	for _, v := range splitStr(s) {
		item := strings.Split(v, "=")
		if len(item) == 0 {
			continue
		} else if len(item) == 1 {
			item = append(item, "")
		}
		switch item[0] {
		case "local_port":
			l.Port = common.GetIntNoErrByStr(item[1])
		case "local_ip":
			l.Ip = item[1]
		case "password":
			l.Password = item[1]
		case "target_addr":
			l.Target = item[1]
		}
	}
	return l
}

func getAllTitle(content string) (arr []string, err error) {
	var re *regexp.Regexp
	re, err = regexp.Compile(`(?m)^\[[^\[\]\r\n]+]`)
	if err != nil {
		return
	}
	arr = re.FindAllString(content, -1)
	m := make(map[string]bool)
	for _, v := range arr {
		if _, ok := m[v]; ok {
			err = errors.New(fmt.Sprintf("Item names %s are not allowed to be duplicated", v))
			return
		}
		m[v] = true
	}
	return
}

func splitStr(s string) (configDataArr []string) {
	if common.IsWindows() {
		configDataArr = strings.Split(s, "\r\n")
	}
	if len(configDataArr) < 3 {
		configDataArr = strings.Split(s, "\n")
	}
	return
}
