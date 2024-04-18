package service

import (
	"errors"
	"xnps/lib/database/models"
)

func (s *DbUtils) CheckUserName(username string) bool {
	if len(username) > 3 {
		var c int64
		s.GDb.Model(models.SystemConfig{}).Where("web_username = ?", username).Limit(1).Count(&c)
		return c > 0
	}
	return false
}

func (s *DbUtils) AddSysConfig(sCOnf *models.SystemConfig) (sysConfig *models.SystemConfig, err error) {

	if _, err2 := s.GetSystemConfig(); err2 != nil {
		s.GDb.Model(models.SystemConfig{}).Create(sCOnf)
		return sCOnf, nil
	} else {
		err = errors.New("already have system config")
	}
	return
}

func (s *DbUtils) UpdateSysConfig(config *models.SystemConfig) {

	s.GDb.Model(models.SystemConfig{}).Updates(config)
}

func (s *DbUtils) GetSystemConfig() (sys models.SystemConfig, err error) {
	if s.GDb.Model(models.SystemConfig{}).First(&sys).RowsAffected < 1 {
		err = errors.New("have no sys config")
	}
	return
}
