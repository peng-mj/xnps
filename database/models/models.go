package models

import (
	"errors"
	"strings"
	"sync"
	"sync/atomic"
)

// 只执行一次，当数据库中无配置信息时
type SystemConfig struct {
	Id                int64  `gorm:"column:id;type:integer;primaryKey" json:"Id"`
	WebHost           string `gorm:"column:web_host;type:text;not null;default:0.0.0.0" json:"WebHost"` //默认服务地址
	WebPort           int64  `gorm:"column:web_port;type:integer;not null;default:8912" json:"WebPort"` //对外服务默认8912
	WebUserName       string `gorm:"column:web_username;type:text;not null;default:admin" json:"WebUserName"`
	WebPassword       string `gorm:"column:web_password;type:text;not null;default:123" json:"WebPassword"`         //因为是sha256加密，所以需要考虑密码重置的情况
	WebOpenCaptcha    bool   `gorm:"column:web_open_captcha;type:integer;not null;default:1" json:"WebOpenCaptcha"` //是否开启验证码校验
	AuthCryptKey      string `gorm:"column:auth_Crypt_key;type:text;not null;default:awdsvthgfd" json:"AuthCryptKey"`
	AllowPorts        string `gorm:"column:allow_ports;type:text;not null;default:8000-10000,20000-30000" json:"AllowPorts"`
	PublicKey         string `gorm:"column:public_key;type:text;not null;default:3d2dw2" json:"PublicKey"`
	BridgeType        string `gorm:"column:bridge_type;type:text;not null;default:tcp" json:"BridgeType"` //tcp、udp、socket、kcp
	BridgePort        int    `gorm:"column:bridge_port;type:integer;not null;default:8913" json:"BridgePort"`
	BridgeHost        string `gorm:"column:bridge_host;type:text;not null;default:0.0.0.0" json:"BridgeHost"` //
	LogLevel          int    `gorm:"column:block_type;type:integer;not null;default:6" json:"LogLevel"`
	LogPath           string `gorm:"column:block_type;type:text;not null;default:" json:"LogPath"`
	MaxClient         int    `gorm:"column:block_type;type:integer;not null;default:1" json:"MaxClient"` //这里可以根据不同性能设备做一下说明
	MaxConn           int    `gorm:"column:block_type;type:integer;not null;default:1" json:"MaxConn"`
	DisConnTimeoutSec int    `gorm:"column:block_type;type:integer;not null;default:1" json:"DisConnTimeoutSec"`
	AllowRegistration bool   `gorm:"column:web_open_captcha;type:integer;not null;default:1" json:"AllowRegistration"` //是否允许用户注册，不允许，仅在无系统配置数据的时候允许
}

func (*SystemConfig) TableName() string {
	return "system_config"
}

func (s *Client) CutConn() {
	atomic.AddUint32(&s.NowConn, 1)

}

func (s *Client) AddConn() {
	atomic.AddUint32(&s.NowConn, -1)
}

func (s *Client) GetConn() bool {
	if s.MaxConn == 0 || s.NowConn < s.MaxConn {
		s.CutConn()
		return true
	}
	return false
}

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
