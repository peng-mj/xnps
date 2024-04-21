package models

type Tunnel struct {
	Id            int64  `gorm:"column:id;type:integer;auto_increment;not null;primaryKey;" json:"id"`
	ClientId      int64  `gorm:"column:client_id;type:integer;not null" json:"clientId"`
	ServerPort    uint16 `gorm:"column:server_port;type:integer;not null;default:8080" json:"ServerPort"`
	Mode          string `gorm:"column:mode;type:text;not null;default:tcp" json:"mode"`
	ConnLimitRate int    `gorm:"column:conn_limit;type:integer;not null;default:60" json:"connLimitRate"`
	Connected     bool   `gorm:"column:status;not null;default:false" json:"connected"`
	ServerPorts   string `gorm:"column:ports;type:text;not null;default:80" json:"ports"`
	TargetIp      string `gorm:"column:target_addr;type:text;not null;default:" json:"target_addr"`
	TargetPort    string `gorm:"column:target_addr;type:text;not null;default:" json:"target_port"`
	Remark        string `gorm:"column:remark;type:text;not null;default:" json:"remark"`
	Valid         bool   `gorm:"column:valid;not null;default:true" json:"valid"`
}

func (*Tunnel) TableName() string {
	return "tunnel"
}
