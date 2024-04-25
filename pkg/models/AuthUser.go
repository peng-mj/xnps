package models

import (
	"time"
)

type AuthUser struct {
	Id           int32     `gorm:"column:id;type:int;primaryKey;autoIncrement" json:"-"`
	Uid          string    `gorm:"column:uuid;type:varchar(64);not null;unique;index" json:"uid,omitempty"`
	Level        int32     `gorm:"column:level;type:int;not null;default:99" json:"level,omitempty"`
	Username     string    `gorm:"column:username;type:varchar(30);not null;default:guest" json:"username"`
	Password     string    `gorm:"column:password;type:varchar(128);not null;default:admin" json:"password,omitempty"`
	Emile        string    `gorm:"column:emil;type:varchar(30);not null;default:empty" json:"emile,omitempty"`
	EmileAuthKey string    `gorm:"column:emil_auth_key;type:varchar(40);not null;default:empty" json:"emile_auth_key"` //SMTP or other
	EmileEnable  bool      `gorm:"column:emil_enable;not null;default:false" json:"emile_enable,omitempty"`
	EmileActive  bool      `gorm:"column:emil_active;not null;default:false" json:"emile_active,omitempty"`
	OTAKeys      string    `gorm:"column:ota_keys;type:varchar(255);not null;default:xnps" json:"ota_keys,omitempty"`
	OTAEnable    bool      `gorm:"column:ota_keys;not null;default:false" json:"ota_enable,omitempty"`
	LastLoginIp  string    `gorm:"column:last_login_ip;type:varchar(20);not null;default:127.0.0.1" json:"last_login_ip,omitempty"`
	CreateAt     time.Time `gorm:"column:create_at;autoCreateTime" json:"create_at"`
	LastLoginAt  int64     `gorm:"column:last_login_at;type:int;not null" json:"last_login_at"`
	ExpirationAt int64     `gorm:"column:expiration_at;type:int;not null" json:"expiration_at"`
	UsagePorts   string    `gorm:"column:usage_ports;type:varchar(512);not null;default:8080-9000" json:"usage_ports,omitempty"`
	MaxConn      int       `gorm:"column:max_conn;type:int;not null;default:100" json:"max_conn"`
	CurConn      int       `gorm:"-" json:"cur_conn"` //not save to database
	Valid        bool      `gorm:"column:valid;not null;default:false" json:"valid"`
}

func (*AuthUser) TableName() string {
	return "auth_user"
}
