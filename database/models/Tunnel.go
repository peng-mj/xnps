package models

import "sync"

type Tunnel struct {
	Id            int64   `gorm:"column:id;type:integer;auto_increment;not null;primaryKey;" json:"id"`
	ClientId      int64   `gorm:"column:client_id;type:integer;not null" json:"clientId"` //所属客户端ID
	Valid         bool    `gorm:"column:valid;not null;default:true" json:"valid"`
	ServerPort    uint16  `gorm:"column:server_port;type:integer;not null;default:8080" json:"ServerPort"`
	ServerIp      string  `gorm:"column:server_ip;type:text;not null;default:0.0.0.0" json:"serverIp"`
	Mode          string  `gorm:"column:mode;type:text;not null;default:" json:"mode"`                     //隧道模式 tcp、udp
	ConnLimitRate int     `gorm:"column:conn_limit;type:integer;not null;default:60" json:"connLimitRate"` //每分钟的连接数量的限制
	Status        bool    `gorm:"column:status;not null;default:false" json:"status"`
	Ports         string  `gorm:"column:ports;type:text;not null;default:80" json:"Ports"`
	Password      string  `gorm:"column:passwd;type:text;not null;default:" json:"Password"` //p2p or secret must use passwd，it must be sha256 not be plaintext password
	Remark        string  `gorm:"column:remark;type:text;not null;default:" json:"Name"`
	TargetAddr    string  `gorm:"column:target_addr;type:text;not null;default:" json:"TargetAddr"`
	NoStore       bool    `gorm:"column:no_store;type:integer;not null;default:0" json:"NoStore"`
	IsHttp        bool    `gorm:"column:is_http;type:integer;not null;default:0" json:"IsHttp"`
	LocalPath     string  `gorm:"column:local_path;type:text;not null;default:" json:"LocalPath"`
	StripPre      string  `gorm:"column:strip_pre;type:text;not null;default:" json:"StripPre"`
	Flow          *Flow   `gorm:"-" json:"-"`
	Client        *Client `gorm:"-" json:"-"`
	Target        *Target `gorm:"-" json:"-"`
	//Health       `gorm:"-" json:"-"`
	sync.RWMutex `gorm:"-" json:"-"`
}

func (*Tunnel) TableName() string {
	return "tunnel"
}
