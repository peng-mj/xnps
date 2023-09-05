package models

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"xnps/lib/database"

	"github.com/pkg/errors"
	"xnps/lib/rate"
)

type Flow struct {
	ExportFlow int64 //出口流浪
	InletFlow  int64 //入口流量
	FlowLimit  int64 //流量限制
	sync.RWMutex
}

func (s *Flow) Add(in, out int64) {
	s.Lock()
	defer s.Unlock()
	s.InletFlow += int64(in)
	s.ExportFlow += int64(out)
}

type Config struct {
	U        string //web的用户名
	P        string //web的密码
	Compress bool
	Crypt    bool
}

type Client2 struct {
	Id              int64    `gorm:"column:id;type:int;auto_increment;not null;primaryKey;" json:"Id"`
	VerifyKey       string   `gorm:"column:verify_key;type:text;not null" json:"VerifyKey"`
	Addr            string   `gorm:"column:addr;type:text" json:"Addr"`
	Remark          string   `gorm:"column:remark;type:text" json:"Remark"`
	Valid           *int     `gorm:"column:valid;type:integer;default:0;not null" json:"Valid"`
	Connected       *int     `gorm:"column:connected;type:integer;default:0;not null" json:"Connected"`
	RateLimit       *int     `gorm:"column:rate_limit;type:integer;default:0;not null" json:"RateLimit"`
	FlowExport      *float32 `gorm:"column:flow_export;type:real;not null;default:0" json:"FlowExport"`
	FlowInle        *float32 `gorm:"column:flow_inle;type:real;not null;default:0" json:"FlowInle"`
	NowRate         *float32 `gorm:"column:now_rate;type:real;default:0;not null" json:"NowRate"`
	MaxConn         int32    `gorm:"column:max_conn;type:integer;default:100;not null" json:"MaxConn"`
	NowConn         int32    `gorm:"column:now_conn;type:integer;default:0;not null" json:"NowConn"`
	WebUser         string   `gorm:"column:web_user;type:text" json:"WebUser"`
	WebPasswd       string   `gorm:"column:web_passwd;type:text" json:"WebPasswd"`
	AllowFileConfig string   `gorm:"column:allow_file_config;type:text" json:"AllowFileConfig"`
	MaxTunnelNum    string   `gorm:"column:max_tunnel_num;type:integer" json:"MaxTunnelNum"`
	Version         string   `gorm:"column:version;type:text" json:"Version"`
	BlackId         *int     `gorm:"column:black_id;type:integer;default:0;not null" json:"BlackId"`
}

func (*Client2) TableName() string {
	return "client"
}

type Client struct {
	Cnf             *Config
	Id              int        //id
	VerifyKey       string     //verify key
	Addr            string     //the ip of client
	Remark          string     //remark
	Status          bool       //is allow connect
	IsConnect       bool       //is the client connect
	RateLimit       int        //rate /kb
	Flow            *Flow      //flow setting
	Rate            *rate.Rate //rate limit
	NoStore         bool       //no store to file
	NoDisplay       bool       //no display on web
	MaxConn         int        //the max connection num of client allow
	NowConn         int32      //the connection num of now
	WebUserName     string     //the username of web login
	WebPassword     string     //the password of web login
	ConfigConnAllow bool       //is allow connected by config file
	MaxTunnelNum    int
	Version         string
	BlackIpList     []string
	sync.RWMutex
}

func NewClient(vKey string, noStore bool, noDisplay bool) *Client {
	return &Client{
		Cnf:       new(Config),
		Id:        0,
		VerifyKey: vKey,
		Addr:      "",
		Remark:    "",
		Status:    true,
		IsConnect: false,
		RateLimit: 0,
		Flow:      new(Flow),
		Rate:      nil,
		NoStore:   noStore,
		RWMutex:   sync.RWMutex{},
		NoDisplay: noDisplay,
	}
}

func (s *Client) CutConn() {
	atomic.AddInt32(&s.NowConn, 1)
}

func (s *Client) AddConn() {
	atomic.AddInt32(&s.NowConn, -1)
}

func (s *Client) GetConn() bool {
	if s.MaxConn == 0 || int(s.NowConn) < s.MaxConn {
		s.CutConn()
		return true
	}
	return false
}

func (s *Client) HasTunnel(t *Tunnel) (exist bool) {
	file.GetDb().JsonDb.Tasks.Range(func(key, value interface{}) bool {
		v := value.(*Tunnel)
		if v.Client.Id == s.Id && v.Port == t.Port && t.Port != 0 {
			exist = true
			return false
		}
		return true
	})
	return
}

func (s *Client) GetTunnelNum() (num int) {
	file.GetDb().JsonDb.Tasks.Range(func(key, value interface{}) bool {
		v := value.(*Tunnel)
		if v.Client.Id == s.Id {
			num++
		}
		return true
	})
	return
}

func (s *Client) HasHost(h *Host) bool {
	var has bool
	file.GetDb().JsonDb.Hosts.Range(func(key, value interface{}) bool {
		v := value.(*Host)
		if v.Client.Id == s.Id && v.Host == h.Host && h.Location == v.Location {
			has = true
			return false
		}
		return true
	})
	return has
}

type Tunnel struct {
	Id           int
	Port         int
	ServerIp     string
	Mode         string
	Status       bool
	RunStatus    bool
	Client       *Client
	Ports        string
	Flow         *Flow
	Password     string
	Remark       string
	TargetAddr   string
	NoStore      bool
	IsHttp       bool
	LocalPath    string
	StripPre     string
	Target       *Target
	MultiAccount *MultiAccount
	Health
	sync.RWMutex
}

type Health struct {
	HealthCheckTimeout  int
	HealthMaxFail       int
	HealthCheckInterval int
	HealthNextTime      time.Time
	HealthMap           map[string]int
	HttpHealthUrl       string
	HealthRemoveArr     []string
	HealthCheckType     string
	HealthCheckTarget   string
	sync.RWMutex
}

type Host struct {
	Id           int
	Host         string //host
	HeaderChange string //header change
	HostChange   string //host change
	Location     string //url router
	Remark       string //remark
	Scheme       string //http https all
	CertFilePath string
	KeyFilePath  string
	NoStore      bool
	IsClose      bool
	Flow         *Flow
	Client       *Client
	Target       *Target //目标
	Health       `json:"-"`
	sync.RWMutex
}

type Target struct {
	nowIndex   int
	TargetStr  string
	TargetArr  []string
	LocalProxy bool
	sync.RWMutex
}

type MultiAccount struct {
	AccountMap map[string]string // multi account and pwd
}

func (s *Target) GetRandomTarget() (string, error) {
	if s.TargetArr == nil {
		s.TargetArr = strings.Split(s.TargetStr, "\n")
	}
	if len(s.TargetArr) == 1 {
		return s.TargetArr[0], nil
	}
	if len(s.TargetArr) == 0 {
		return "", errors.New("all inward-bending targets are offline")
	}
	s.Lock()
	defer s.Unlock()
	if s.nowIndex >= len(s.TargetArr)-1 {
		s.nowIndex = -1
	}
	s.nowIndex++
	return s.TargetArr[s.nowIndex], nil
}
