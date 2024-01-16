package models

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
