package service

import (
	"errors"
	"tunpx/pkg/crypt"
	"tunpx/pkg/models"
	"tunpx/web/dto"
)

type AuthUser struct {
	*Base
}

func NewAuthUser(db *Base) *AuthUser {
	a := &AuthUser{}
	a.Base = db
	return a
}

func (s *AuthUser) CheckPasswd(auth *dto.LoginReq) (user *models.AuthUser, code dto.RspCode) {
	user = new(models.AuthUser)
	if s.Orm(models.AuthUser{}).Where("username = ? or (emil_enable = 1 and emil = ?) ", auth.Username, auth.Password).First(&user).RowsAffected == 0 {
		return nil, dto.RspCode(dto.ErrNotFound)
	}
	var otpCOde, passwd string
	if user.OTAEnable {
		if len(auth.OtpCode) != 6 {
			// return errors.New("need otp code")
			return nil, dto.RspCode(dto.NeedOtpCode)

		}
		// TODO: design a otp login
		//	get otp code
		// otpCOde={otp code}
		passwd = crypt.Sha256(auth.Username + "." + user.Password + "." + otpCOde)
	} else {
		passwd = crypt.Sha256(auth.Username + "." + user.Password)
	}
	if auth.Password != passwd {
		return nil, dto.RspCode(dto.NeedOtpCode)
	}
	return user, dto.RspCode(200)
}
func (s *AuthUser) GetUserByUid(uid string) (user models.AuthUser, err error) {
	if s.Orm(models.AuthUser{}).Where("uid = ?", uid).First(&user).RowsAffected == 0 {
		err = errors.New("user not found")
	}
	return
}

// GetAllUser just admin
func (s *AuthUser) GetAllUser() (user []models.AuthUser) {
	user = make([]models.AuthUser, 0)
	s.Orm(models.AuthUser{}).Find(&user)
	return
}
func (s *AuthUser) Create(user *models.AuthUser) error {
	err := s.Orm(models.AuthUser{}).Create(user).Error
	return err
}
