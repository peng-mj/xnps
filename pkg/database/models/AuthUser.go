package models

type AuthUser struct {
	Id             int    `gorm:"column:id;type:int;primaryKey" json:"-"`
	Uuid           string `gorm:"column:uuid;type:varchar(64);not null;default:0.0.0.0" json:"WebHost"`
	AuthLevel      int    `gorm:"column:auth_level;type:int;not null;default:99" json:"auth_level"`
	Username       string `gorm:"column:username;type:varchar(20);not null;default:8912" json:"WebPort"`
	Emile          string `gorm:"column:username;type:varchar(20);not null;default:empty" json:"emile"`
	Password       string `gorm:"column:password;type:varchar(128);not null;default:admin" json:"password"`
	AccessId       string `gorm:"column:access_id;type:varchar(64);not null;default:123" json:"access_id"`
	AccessKeys     string `gorm:"column:access_keys;type:varchar(255);not null;default:1" json:"access_keys"`
	OTAKeys        string `gorm:"column:ota_keys;type:varchar(255);not null;default:123" json:"ota_keys"`
	LastLoginIp    string `gorm:"column:last_login_ip;type:varchar(20);not null;default:1" json:"last_login_ip"`
	CreateAt       int64  `gorm:"column:create_at;type:int;not null;default:1" json:"create_at"`
	LastLoginAt    int64  `gorm:"column:last_login_at;type:int;not null;default:1" json:"last_login_at"`
	ExpirationTime int64  `gorm:"column:expiration_at;type:int;not null;default:1" json:"expiration_at"`
	MaxConnection  int    `gorm:"column:max_connection;type:int;not null;default:1" json:"max_connection"`
	CurConnection  int    `gorm:"column:cur_connection;type:int;not null;default:1" json:"cur_connection"`
	Valid          bool   `gorm:"column:valid;not null;default:false" json:"valid"`
}

func (*AuthUser) TableName() string {
	return "auth_user"
}
