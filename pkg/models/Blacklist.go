package models

type Blacklist struct {
	Id         int64  `gorm:"column:id;type:integer;auto_increment;not null;primaryKey;" json:"id"`
	BlockType  int64  `gorm:"column:block_type;type:integer;not null;default:1" json:"block_type"`
	FirewallId int64  `gorm:"column:firewall_id;type:integer;not null;default:1" json:"firewall_id"`
	SourceIp   string `gorm:"column:source_ip;type:text;not null;default:1" json:"source_ip"`
	TargetIp   string `gorm:"column:target_ip;type:text;not null;default:1" json:"target_ip"`
	Location   int64  `gorm:"column:location;type:text;not null;default:1" json:"location"`
	Belong     int64  `gorm:"column:belong;type:integer;not null;default:1" json:"belong"`
	CreateAt   int64  `gorm:"column:create_at;type:integer;not null;default:1" json:"create_at"`
}

func (*Blacklist) TableName() string {
	return "blacklist"
}
