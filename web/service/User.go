package service

import (
	"errors"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"os"
	"xnps/lib/database"
	"xnps/lib/database/models"
)

type DbUtils struct {
	GDb *gorm.DB
}

type User Base

// init csv from file
func GetDb() *DbUtils {
	if database.Db == nil {
		//logs.Info("数据库未打开")
		slog.Error("数据库未打开")

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
