package service

import (
	"fmt"
	"tunpx/pkg/models"
)

type System struct {
	Base
}

func NewSystem(db *Base) *System {
	c := &System{}
	c.Service(db.Driver)
	return c
}
func (s *System) CreateInit(conf *models.Config) error {
	var c int64
	s.Orm(models.Config{}).Count(&c)
	if c != 0 {
		return fmt.Errorf("system aready have init config")
	}
	return s.Orm(models.Config{}).Create(conf).Error
}

func (s *System) Update(conf *models.Config) error {
	return s.Orm(models.Config{}).Updates(conf).Error
}
func (s *System) Get() (conf *models.Config) {
	s.Orm(models.Config{}).First(conf)
	return
}
func (s *System) Check() (c int64) {
	s.Orm(models.Config{}).Limit(1).Count(&c)
	return
}
