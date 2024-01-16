package Mapper

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"gorm.io/gorm"
	"os"
	"xnps/database"
	"xnps/database/models"
)

type DbUtils struct {
	GDb *gorm.DB
	//JsonDb *JsonDb
}

// init csv from file
func GetDb() *DbUtils {
	if database.Db == nil {
		logs.Info("数据库未打开")

		os.Exit(-1)
	}
	return database.Db
}

func (s *DbUtils) GetPasswdByUser(user string) (passwd string, err error) {
	sys := new(models.SystemConfig)
	if s.GDb.Model(models.SystemConfig{}).Where("web_username = ?", user).First(sys).RowsAffected > 0 {
		return sys.WebPassword, nil
	} else {
		return "", errors.New("have no user named " + user)
	}
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
