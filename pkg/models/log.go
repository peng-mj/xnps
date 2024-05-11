package models

import (
	"gorm.io/gorm"
	"time"
)

type Log struct {
	Id       int32     `gorm:"column:id;type:int;primaryKey;autoIncrement" json:"-"`
	App      string    `gorm:"column:app;type:varchar(64);not null;unique;index" json:"app,omitempty"`
	Type     string    `gorm:"column:type;type:varchar(64);not null;unique;index" json:"type,omitempty"`
	Code     string    `gorm:"column:code;type:varchar(64);not null;unique;index" json:"code,omitempty"`
	Msg      string    `gorm:"column:msg;type:varchar(64);not null;unique;index" json:"msg,omitempty"`
	CreateAt time.Time `gorm:"column:create_at;autoCreateTime" json:"create_at,omitempty"`
}

func (*Log) TableName() string {
	return "log"
}

func (l *Log) BeforeCreate(*gorm.DB) error {
	l.CreateAt = time.Now()
	return nil
}
