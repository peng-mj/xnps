package models

import (
	"gorm.io/gorm"
	"time"
)

type Config struct {
	Id       int32     `gorm:"column:id;type:int;primaryKey;autoIncrement" json:"-"`
	OrgName  string    `gorm:"column:org_name;type:varchar(128);not null" json:"org_name"`
	CreateAt time.Time `gorm:"column:create_at;autoCreateTime" json:"create_at"`
	// AppKey     string    `gorm:"column:app_key;type:varchar(255);not null" json:"app_key"` // Used to sign the token
	BridgePort int    `gorm:"column:bridge_port;type:int;not null;default:9871" json:"bridge_port"`
	UsagePorts string `gorm:"column:usage_port;type:varchar(512);not null;default:8000-8100" json:"usage_ports"`
	WebPort    int    `gorm:"column:web_port;type:int;not null;default:9870" json:"web_port"`
	Remark     string `gorm:"column:remark;type:varchar(255);" json:"remark"`
}

func (*Config) TableName() string {
	return "config"
}

func (m *Config) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateAt = time.Now()
	return
}
