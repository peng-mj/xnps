package models

import (
	"gorm.io/gorm"
	"time"
)

type Group struct {
	Id           int64     `gorm:"column:id;type:int;auto_increment;not null;primaryKey;" json:"id"`
	Name         string    `gorm:"column:name;type:varchar(40);not null;" json:"name"`
	Uid          int64     `gorm:"column:uid;type:int;not null;" json:"uid"` // user's uuid
	UsagePorts   string    `gorm:"column:ports;type:text;not null;" json:"ports"`
	GroupType    string    `gorm:"column:group_type;type:varchar(10);not null;default:dft" json:"group_type"`
	CreateAt     time.Time `gorm:"column:create_at;autoCreateTime" json:"create_at"`
	ModifyAt     time.Time `gorm:"column:modify_at;autoUpdateTime" json:"modify_at"`
	MaxClientNum int32     `gorm:"column:max_client_num;type:int;not null;default:0" json:"max_client_num"`
	CurClientNum int32     `gorm:"column:cur_client_num;type:int;not null;default:0" json:"cur_client_num"` // 可能不需要写入数据库
	OnClientNum  int32     `gorm:"column:on_client_num;type:int;not null;default:0" json:"on_client_num"`   // 可能不需要写入数据库
	Valid        bool      `gorm:"column:valid;not null;default:true" json:"valid"`
	Remark       string    `gorm:"column:remark;type:text;not null;" json:"remark"`
}

func (s *Group) TableName() string {

	return "client_group"
}

func (s *Group) BeforeCreate(tx *gorm.DB) (err error) {
	s.CreateAt = time.Now()
	return
}

func (s *Group) BeforeUpdate(tx *gorm.DB) (err error) {
	s.ModifyAt = time.Now()
	return
}
