package service

import (
	"errors"
	"xnps/pkg/models"
)

type AuthUser struct {
	*Base
}

func NewAuthUser(db *Base) *AuthUser {
	a := &AuthUser{}
	a.Base = db
	return a
}

func (s *AuthUser) CheckPasswd(name, password, otp string) error {
	var user models.AuthUser
	if s.Orm(models.AuthUser{}).Where("username = ? or (emil_enable = 1 and emil = ?) ", name, name).First(&user).RowsAffected == 0 {
		return errors.New("user not found")
	}
	if password != user.Password {
		return errors.New("password error")
	}

	return nil
}
func (s *AuthUser) GetUserByUid(uid string) (user models.AuthUser, err error) {
	if s.Orm(models.AuthUser{}).Where("uid = ?", uid).First(&user).RowsAffected == 0 {
		err = errors.New("user not found")
	}
	return
}

// GetAllUser just admin
func (s *AuthUser) GetAllUser() (user []models.AuthUser, err error) {
	user = make([]models.AuthUser, 0)
	if s.Orm(models.AuthUser{}).Find(&user).RowsAffected == 0 {
		err = errors.New("user not found")
	}
	return
}
