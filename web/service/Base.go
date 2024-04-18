package service

import "gorm.io/gorm"

type Base struct {
	GDb *gorm.DB
}
