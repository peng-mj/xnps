package models

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
