package models

import (
	"github.com/pkg/errors"
	"strings"
	"sync"
	"sync/atomic"
	"xnps/lib/rate"
)

type Flow struct {
	Id         int64 `gorm:"column:id;type:int;auto_increment;not null;primaryKey;" json:"Id"`
	ClientId   int64 `gorm:"column:export;type:int;auto_increment;not null;primaryKey;" json:"ClientId"`
	ExportFlow int64 `gorm:"column:export;type:int;auto_increment;not null;primaryKey;" json:"ExportFlow"` //出口流浪
	InletFlow  int64 `gorm:"column:inlet;type:int;auto_increment;not null;primaryKey;" json:"InletFlow"`   //入口流量
	FlowLimit  int64 `gorm:"column:limit;type:int;auto_increment;not null;primaryKey;" json:"FlowLimit"`   //流量限制
	sync.RWMutex
}

func (s *Flow) Add(in, out int64) {
	s.Lock()
	defer s.Unlock()
	s.InletFlow += in
	s.ExportFlow += out
}
func (s *Flow) TableName() string {
	return "flow_info"
}

type Client struct {
	Id                 int64      `gorm:"column:id;type:int;auto_increment;not null;primaryKey;" json:"Id"`
	VerifyKey          string     `gorm:"column:verify_key;type:text;not null" json:"VerifyKey"`
	RemoteAddr         string     `gorm:"column:remote_addr;type:text;not null;default: " json:"RemoteAddr"` //客户端真实IP
	Remark             string     `gorm:"column:remark;type:text;not null;default: " json:"Remark"`
	Valid              bool       `gorm:"column:valid;type:integer;default:0;not null" json:"Valid"`
	Connected          bool       `gorm:"column:connected;type:integer;default:0;not null" json:"Connected"`
	Crypt              bool       `gorm:"column:crypt;type:integer;not null;default:" json:"Crypt"`       //是否加密
	Compress           bool       `gorm:"column:compress;type:integer;not null;default:" json:"Compress"` //是否压缩
	RateLimit          int        `gorm:"column:rate_limit;type:integer;default:0;not null" json:"RateLimit"`
	FlowExport         float32    `gorm:"column:flow_export;type:real;not null;default:0" json:"FlowExport"`
	FlowInle           float32    `gorm:"column:flow_inle;type:real;not null;default:0" json:"FlowInle"`
	NowRate            float32    `gorm:"column:now_rate;type:real;default:0;not null" json:"NowRate"`
	MaxConn            int        `gorm:"column:max_conn;type:integer;default:100;not null" json:"MaxConn"`
	NowConn            int32      `gorm:"column:now_conn;type:integer;default:0;not null" json:"NowConn"`
	HttpUser           string     `gorm:"column:http_user;type:text;not null;default:user" json:"HttpUser"` //这个用于用户登录
	HttpPasswd         string     `gorm:"column:http_passwd;type:text;not null;default:123" json:"HttpPasswd"`
	AllowUseConfigFile bool       `gorm:"column:allow_file_config;type:integer;default:1;not null" json:"AllowUseConfigFile"`
	MaxTunnelNum       int        `gorm:"column:max_tunnel_num;type:integer;default:100;not null" json:"MaxTunnelNum"`
	Version            string     `gorm:"column:version;type:text;default:Null;not null" json:"Version"`
	BlackId            int        `gorm:"column:black_id;type:integer;default:0;not null" json:"BlackId"`
	ActiveTime         int64      `gorm:"column:active_time;type:integer;default:1672502400;not null" json:"ActiveTime"`
	Flow               *Flow      `gorm:"-" json:"-"`
	Rate               *rate.Rate `gorm:"-" json:"-"`
	sync.RWMutex
}

func (*Client) TableName() string {
	return "client"
}

type UserInfo struct {
	Id            int64  `gorm:"column:id;type:integer;primaryKey" json:"Id"`
	Valid         bool   `gorm:"column:valid;type:integer;not null;default:1" json:"Valid"`
	UserName      string `gorm:"column:username;type:text;not null;default:1" json:"UserName"`
	Passwd        string `gorm:"column:passwd;type:text;not null;default:1" json:"Passwd"` //sha256加密
	CreateTime    int64  `gorm:"column:create_time;type:integer;not null;default:1" json:"CreateTime"`
	LastLoginTime int64  `gorm:"column:last_login_time;type:integer;not null;default:1" json:"lastLoginTime"`
	LastLoginIp   string `gorm:"column:last_login_ip;type:text;not null;default:1" json:"LastLoginIp"`
	AuthType      string `gorm:"column:auth_type;type:text;not null;default:1" json:"AuthType"`
}

func (*UserInfo) TableName() string {
	return "user_info"
}

type Tunnel struct {
	Id              int64         `gorm:"column:id;type:integer;primaryKey" json:"Id"`
	Valid           bool          `gorm:"column:valid;type:integer;not null;default:1" json:"Valid"`
	ClientId        int64         `gorm:"column:client_id;type:integer;not null" json:"ClientId"`
	ServerPort      int           `gorm:"column:server_port;type:integer;not null;default:8080" json:"ServerPort"`
	ServerIp        string        `gorm:"column:server_ip;type:text;not null;default:" json:"ServerIp"`
	Mode            string        `gorm:"column:mode;type:text;not null;default:" json:"Mode"`
	ConnLimitPerMin int           `gorm:"column:conn_limit;type:integer;not null;default:60" json:"ConnLimitPerMin"` //每分钟的连接数量的限制
	Status          bool          `gorm:"column:status;type:integer;not null;default:" json:"Status"`
	RunStatus       bool          `gorm:"column:run_status;type:integer;not null;default:" json:"RunStatus"` //运行状态
	Ports           string        `gorm:"column:ports;type:text;not null;default:80" json:"Ports"`           //仅适用于p2p和私密代理
	Password        string        `gorm:"column:passwd;type:text;not null;default:" json:"Password"`         //p2p or secret must use passwd，it must be sha256 not be plaintext password
	Remark          string        `gorm:"column:remark;type:text;not null;default:" json:"Remark"`
	TargetAddr      string        `gorm:"column:target_addr;type:text;not null;default:" json:"TargetAddr"`
	NoStore         bool          `gorm:"column:no_store;type:integer;not null;default:0" json:"NoStore"`
	IsHttp          bool          `gorm:"column:is_http;type:integer;not null;default:0" json:"IsHttp"`
	LocalPath       string        `gorm:"column:local_path;type:text;not null;default:" json:"LocalPath"`
	StripPre        string        `gorm:"column:strip_pre;type:text;not null;default:" json:"StripPre"`
	Flow            *Flow         `gorm:"-" json:"-"`
	Client          *Client       `gorm:"ForeignKey:client_id" json:"-"`
	Target          *Target       `gorm:"-" json:"-"`
	MultiAccount    *MultiAccount `gorm:"-" json:"-"`
	//Health       `gorm:"-" json:"-"`
	sync.RWMutex `gorm:"-" json:"-"`
}

func (*Tunnel) TableName() string {
	return "tunnel"
}

type Firewall struct {
	Id            int64  `gorm:"column:id;type:integer;primaryKey" json:"Id"`
	Valid         bool   `gorm:"column:valid;type:integer;not null;default:1" json:"Valid"`
	UpdateTime    int64  `gorm:"column:valid;type:integer;not null;default:1" json:"UpdateTime"`
	FType         string `gorm:"column:type;type:text;not null;default:1" json:"FType"` //防火墙类型，白名单还是黑名单
	IpRules       string `gorm:"column:ip_rules;type:text;not null;default:1" json:"IpRules"`
	LocationRules string `gorm:"column:location_rules;type:text;not null;default:0" json:"LocationRules"`
}

func (*Firewall) TableName() string {
	return "firewall"
}

type BlockListInfo struct {
	Id         int64  `gorm:"column:id;type:integer;primaryKey" json:"Id"`
	BlockType  int64  `gorm:"column:block_type;type:integer;not null;default:1" json:"BlockType"`
	SourceIp   string `gorm:"column:ip_info;type:text;not null;default:1" json:"SourceIp"`
	TargetIp   string `gorm:"column:ip_info;type:text;not null;default:1" json:"TargetIp"`
	Location   int64  `gorm:"column:location;type:text;not null;default:1" json:"Location"`
	Belong     int64  `gorm:"column:belong;type:integer;not null;default:1" json:"Belong"`
	CreateTime int64  `gorm:"column:create_time;type:integer;not null;default:1" json:"UpdateTime"`
}

func (*BlockListInfo) TableName() string {
	return "block_recode"
}

// 只执行一次，当数据库中无配置信息时
type SystemConfig struct {
	Id                int64  `gorm:"column:id;type:integer;primaryKey" json:"Id"`
	WebHost           int64  `gorm:"column:block_type;type:integer;not null;default:1" json:"WebHost"`
	WebPort           int64  `gorm:"column:block_type;type:integer;not null;default:1" json:"WebPort"`
	WebUserName       int64  `gorm:"column:block_type;type:integer;not null;default:1" json:"WebUserName"`
	WebPassword       string `gorm:"column:block_type;type:integer;not null;default:1" json:"WebPassword"`  //因为是sha256加密，所以需要考虑密码重置的情况
	WebIsCaptcha      bool   `gorm:"column:block_type;type:integer;not null;default:1" json:"WebIsCaptcha"` //是否开启验证码校验
	AuthCryptKey      string `gorm:"column:block_type;type:integer;not null;default:1" json:"AuthCryptKey"`
	AllowPorts        string `gorm:"column:block_type;type:integer;not null;default:1" json:"AllowPorts"`
	PublicKey         string `gorm:"column:block_type;type:integer;not null;default:1" json:"PublicKey"`
	BridgeType        string `gorm:"column:block_type;type:integer;not null;default:1" json:"BlockType"`
	BridgePort        string `gorm:"column:block_type;type:integer;not null;default:1" json:"BridgePort"`
	BridgeHost        string `gorm:"column:block_type;type:integer;not null;default:1" json:"BridgeHost"`
	LogLevel          string `gorm:"column:block_type;type:integer;not null;default:1" json:"LogLevel"`
	LogPath           string `gorm:"column:block_type;type:integer;not null;default:1" json:"LogPath"`
	MaxClient         int64  `gorm:"column:block_type;type:integer;not null;default:1" json:"MaxClient"`
	MaxConn           int64  `gorm:"column:block_type;type:integer;not null;default:1" json:"MaxConn"`
	DisConnTimeoutSec int64  `gorm:"column:block_type;type:integer;not null;default:1" json:"DisConnTimeoutSec"`

	//BlockType  int64  `gorm:"column:block_type;type:integer;not null;default:1" json:"BlockType"`
	//BlockType  int64  `gorm:"column:block_type;type:integer;not null;default:1" json:"BlockType"`
}

func (*SystemConfig) TableName() string {
	return "system_config"
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

//
//func (s *Client) HasTunnel(t *Tunnel) (exist bool) {
//
//	database.GetDb().JsonDb.Tasks.Range(func(key, value interface{}) bool {
//		v := value.(*Tunnel)
//		if v.Client.Id == s.Id && v.ServerPort == t.ServerPort && t.ServerPort != 0 {
//			exist = true
//			return false
//		}
//		return true
//	})
//	return
//}
//
//// 获取隧道数量
//func (s *Client) GetTunnelNum() (num int) {
//	database.GetDb().JsonDb.Tasks.Range(func(key, value interface{}) bool {
//		v := value.(*Tunnel)
//		if v.Client.Id == s.Id {
//			num++
//		}
//		return true
//	})
//	return
//}

//
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
//
//type Host struct {
//	Id           int
//	Host         string //host
//	HeaderChange string //header change
//	HostChange   string //host change
//	Location     string //url router
//	Remark       string //remark
//	Scheme       string //http https all
//	CertFilePath string
//	KeyFilePath  string
//	NoStore      bool
//	IsClose      bool
//	Flow         *Flow
//	Client       *Client
//	Target       *Target //目标
//	//Health       `json:"-"`
//	sync.RWMutex
//}

// 这个是用来生成 可用的端口的地址
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
