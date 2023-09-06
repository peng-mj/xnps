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
	User     string //web的用户名
	Passwd   string //web的密码
	Compress bool
	Crypt    bool
}

type Client struct {
	Id                 int64      `gorm:"column:id;type:int;auto_increment;not null;primaryKey;" json:"Id"`
	VerifyKey          string     `gorm:"column:verify_key;type:text;not null" json:"VerifyKey"`
	Addr               string     `gorm:"column:addr;type:text;not null;default: " json:"Addr"`
	Remark             string     `gorm:"column:remark;type:text;not null;default: " json:"Remark"`
	Valid              bool       `gorm:"column:valid;type:integer;default:0;not null" json:"Valid"`
	Connected          bool       `gorm:"column:connected;type:integer;default:0;not null" json:"Connected"`
	RateLimit          int        `gorm:"column:rate_limit;type:integer;default:0;not null" json:"RateLimit"`
	FlowExport         float32    `gorm:"column:flow_export;type:real;not null;default:0" json:"FlowExport"`
	FlowInle           float32    `gorm:"column:flow_inle;type:real;not null;default:0" json:"FlowInle"`
	NowRate            float32    `gorm:"column:now_rate;type:real;default:0;not null" json:"NowRate"`
	MaxConn            int32      `gorm:"column:max_conn;type:integer;default:100;not null" json:"MaxConn"`
	NowConn            int32      `gorm:"column:now_conn;type:integer;default:0;not null" json:"NowConn"`
	WebUser            string     `gorm:"column:web_user;type:text;not null;default:user" json:"WebUser"`
	WebPasswd          string     `gorm:"column:web_passwd;type:text;not null;default:123" json:"WebPasswd"`
	AllowUseConfigFile bool       `gorm:"column:allow_file_config;type:integer;default:1;not null" json:"AllowUseConfigFile"`
	MaxTunnelNum       int        `gorm:"column:max_tunnel_num;type:integer;default:100;not null" json:"MaxTunnelNum"`
	Version            string     `gorm:"column:version;type:text;default:Null;not null" json:"Version"`
	BlackId            int        `gorm:"column:black_id;type:integer;default:0;not null" json:"BlackId"`
	Flow               *Flow      `gorm:"-" json:"-"`
	Rate               *rate.Rate `gorm:"-" json:"-"`
	sync.RWMutex
}

func (*Client) TableName() string {
	return "client"
}

func (s *Client) CutConn() {
	atomic.AddInt32(&s.NowConn, 1)
}

func (s *Client) AddConn() {
	atomic.AddInt32(&s.NowConn, -1)
}

func (s *Client) GetConn() bool {
	if s.MaxConn == 0 || s.NowConn < s.MaxConn {
		s.CutConn()
		return true
	}
	return false
}

func (s *Client) HasTunnel(t *Tunnel) (exist bool) {

	database.GetDb().JsonDb.Tasks.Range(func(key, value interface{}) bool {
		v := value.(*Tunnel)
		if v.Client.Id == s.Id && v.Port == t.Port && t.Port != 0 {
			exist = true
			return false
		}
		return true
	})
	return
}

// 获取隧道数量
func (s *Client) GetTunnelNum() (num int) {
	database.GetDb().JsonDb.Tasks.Range(func(key, value interface{}) bool {
		v := value.(*Tunnel)
		if v.Client.Id == s.Id {
			num++
		}
		return true
	})
	return
}

type Tunnel struct {
	Id           int64         `gorm:"column:id;type:integer;primaryKey" json:"Id"`
	ClientId     int           `gorm:"column:client_id;type:integer;" json:"ClientId"`
	Port         int           `gorm:"column:port;type:integer;not null;default:8080" json:"Port"`
	ServerIp     string        `gorm:"column:server_ip;type:text;not null;default:" json:"ServerIp"`
	Mode         string        `gorm:"column:mode;type:text;not null;default:" json:"Mode"`
	Status       bool          `gorm:"column:status;type:integer;not null;default:" json:"Status"`
	RunStatus    bool          `gorm:"column:id;type:integer;not null;default:" json:"RunStatus"`
	Client       *Client       `gorm:"-" json:"-"`
	Ports        string        `gorm:"column:id;type:text;not null;default:80" json:"Ports"`
	Flow         *Flow         `gorm:"-" json:"-"`
	Password     string        `gorm:"column:passwd;type:text;not null;default:" json:"Password"` //p2p or secret must use passwd
	Remark       string        `gorm:"column:remark;type:text;not null;default:" json:"Remark"`
	TargetAddr   string        `gorm:"column:target_addr;type:text;not null;default:" json:"TargetAddr"`
	NoStore      bool          `gorm:"column:no_store;type:integer;not null;default:0" json:"NoStore"`
	IsHttp       bool          `gorm:"column:is_http;type:integer;not null;default:0" json:"IsHttp"`
	LocalPath    string        `gorm:"column:local_path;type:text;not null;default:" json:"LocalPath"`
	StripPre     string        `gorm:"column:strip_pre;type:text;not null;default:" json:"StripPre"`
	Target       *Target       `gorm:"-" json:"-"`
	MultiAccount *MultiAccount `gorm:"-" json:"-"`
	//Health       `gorm:"-" json:"-"`
	sync.RWMutex `gorm:"-" json:"-"`
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
	//Health       `json:"-"`
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
