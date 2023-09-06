package file

import (
	"github.com/pkg/errors"
	"strings"
	"sync"
	"sync/atomic"
	"xnps/lib/rate"
)

type Flow struct {
	ExportFlow int64
	InletFlow  int64
	FlowLimit  int64
	sync.RWMutex
}

func (s *Flow) Add(in, out int64) {
	s.Lock()
	defer s.Unlock()
	s.InletFlow += int64(in)
	s.ExportFlow += int64(out)
}

type Config struct {
	User     string
	Passwd   string
	Compress bool
	Crypt    bool
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
	GetDb().JsonDb.Tasks.Range(func(key, value interface{}) bool {
		v := value.(*Tunnel)
		if v.Client.Id == s.Id && v.Port == t.Port && t.Port != 0 {
			exist = true
			return false
		}
		return true
	})
	return
}

// 总获取隧道数量
func (s *Client) GetTunnelNum() (num int) {
	GetDb().JsonDb.Tasks.Range(func(key, value interface{}) bool {
		v := value.(*Tunnel)
		if v.Client.Id == s.Id {
			num++
		}
		return true
	})
	return
}

//type Tunnel struct {
//	Id           int64         `gorm:"column:primaryKey;id" json:"Id"`
//	Port         int32         `gorm:"column:port;type:integer;not null;default:8080" json:"Port"`
//	ServerIp     string        `gorm:"column:server_ip;type:integer;not null;default:" json:"ServerIp"`
//	Mode         string        `gorm:"column:mode;type:integer;not null;default:" json:"Mode"`
//	Status       bool          `gorm:"column:status;type:integer;not null;default:" json:"Status"`
//	RunStatus    bool          `gorm:"column:run_status;type:integer;not null;default:" json:"RunStatus"`
//	ClientId     int           `gorm:"column:client_id;type:integer;not null;default:" json:"Client"`
//	Ports        string        `gorm:"column:ports;type:integer;not null;default:80" json:"Ports"`
//	FlowId       int           `gorm:"column:flow_id;type:integer;not null;default:" json:"FlowId"`
//	Password     string        `gorm:"column:passwd;type:integer;not null;default:" json:"Password"`
//	Remark       string        `gorm:"column:remark;type:integer;not null;default:" json:"Remark"`
//	TargetAddr   string        `gorm:"column:targetAddr;type:integer;not null;default:" json:"TargetAddr"`
//	NoStore      bool          `gorm:"column:no_store;type:integer;not null;default:" json:"NoStore"`
//	IsHttp       bool          `gorm:"column:is_http;type:integer;not null;default:" json:"IsHttp"`
//	LocalPath    string        `gorm:"column:local_path;type:integer;not null;default:" json:"LocalPath"`
//	StripPre     string        `gorm:"column:strip_pre;type:integer;not null;default:" json:"StripPre"`
//	Target       *Target       `gorm:"-" json:"Target"`
//	MultiAccount *MultiAccount `gorm:"-" json:"MultiAccount"`
//	Health       `gorm:"-" json:"-"`
//	sync.RWMutex `gorm:"-" json:"-"`
//}

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
	//Health
	sync.RWMutex
}

//type Health struct {
//	HealthCheckTimeout  int
//	HealthMaxFail       int
//	HealthCheckInterval int
//	HealthNextTime      time.Time
//	HealthMap           map[string]int
//	HttpHealthUrl       string
//	HealthRemoveArr     []string
//	HealthCheckType     string
//	HealthCheckTarget   string
//	sync.RWMutex
//}

type Target struct {
	nowIndex   int
	TargetStr  string
	TargetArr  []string
	LocalProxy bool
	sync.RWMutex
}

// 这个是当使用socket5代理模式时，多账户
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
