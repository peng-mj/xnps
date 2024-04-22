package models

type AuthUser struct {
	Id           int32  `gorm:"column:id;type:int;primaryKey;autoIncrement" json:"-"`
	Uid          string `gorm:"column:uuid;type:varchar(64);not null;unique;index" json:"uid"`
	Level        int32  `gorm:"column:level;type:int;not null;default:99" json:"level"`
	Username     string `gorm:"column:username;type:varchar(20);not null;default:8912" json:"WebPort"`
	Emile        string `gorm:"column:emil;type:varchar(20);not null;default:empty" json:"emile"`
	EmileEnable  bool   `gorm:"column:emil_enable;not null;default:false" json:"emile_enable"`
	Password     string `gorm:"column:password;type:varchar(128);not null;default:admin" json:"password"`
	OTAKeys      string `gorm:"column:ota_keys;type:varchar(255);not null;default:xnps" json:"ota_keys"`
	OTAEnable    bool   `gorm:"column:ota_keys;not null;default:false" json:"ota_enable"`
	LastLoginIp  string `gorm:"column:last_login_ip;type:varchar(20);not null;default:1" json:"last_login_ip"`
	CreateAt     int64  `gorm:"column:create_at;type:int;not null;default:1" json:"create_at"`
	LastLoginAt  int64  `gorm:"column:last_login_at;type:int;not null;default:1" json:"last_login_at"`
	ExpirationAt int64  `gorm:"column:expiration_at;type:int;not null;default:1" json:"expiration_at"`
	MaxConn      int    `gorm:"column:max_conn;type:int;not null;default:1" json:"max_conn"`
	CurConn      int    `gorm:"column:cur_conn;type:int;not null;default:1" json:"cur_conn"`
	Valid        bool   `gorm:"column:valid;not null;default:false" json:"valid"`
}

func (*AuthUser) TableName() string {
	return "auth_user"
}
