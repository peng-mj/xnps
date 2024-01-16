package Mapper

import "xnps/database/models"

func (s *DbUtils) CheckUserName(username string) bool {
	if len(username) > 3 {
		c := int64(0)
		s.GDb.Model(models.SystemConfig{}).Where("web_username = ?", username).Limit(1).Count(&c)
		return c > 0
	}
	return false
}
