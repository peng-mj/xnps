package service

import (
	"errors"
	"fmt"
	"xnps/lib/common"
	"xnps/lib/database/models"
)

func (s *DbUtils) GetTunnelListByClientIdWithPage(start, recordsCount int, mode string, clientId int64) ([]models.Tunnel, int) {
	var cli []models.Tunnel
	db := s.GDb.Model(models.Tunnel{}).Where("client_id = ? and valid = 1", clientId)
	if common.IsTunnelMode(mode) {
		db = db.Where("mode = ?", mode)
	}
	db.Offset(start).Limit(recordsCount).Find(&cli)
	return cli, len(cli)
}
func (s *DbUtils) GetAllTunnelList(status int) ([]models.Tunnel, int) {
	var cli []models.Tunnel
	db := s.GDb.Model(models.Tunnel{})
	if status == 0 || status == 1 {
		db = db.Where("valid = ?", status)
	}
	db.Find(&cli)
	return cli, len(cli)
}
func (s *DbUtils) GetAllTunnelNumById(id int64) uint32 {
	var num int64
	s.GDb.Model(models.Tunnel{}).Where("id = ?", id).Count(&num)
	return uint32(num)
}

func (s *DbUtils) NewTunnel(tunnel *models.Tunnel) (err error) {
	if tunnel.ClientId <= 0 {
		return fmt.Errorf("所属客户端不存在，请重新提交")
	}
	if err != nil {
		return
	}
	s.GDb.Model(models.Tunnel{}).Create(&tunnel)
	return
}
func (s *DbUtils) UpdateTunnel(tunnel *models.Tunnel) error {
	if tunnel.Id <= 0 {
		return fmt.Errorf("tunnel id 填写错误，请检查")
	} else if tunnel.ClientId <= 0 {
		return fmt.Errorf("所属客户端不存在，请重新提交")
	}
	s.GDb.Model(models.Tunnel{}).Where("id = ?", tunnel.Id).Omit().Updates(tunnel)
	return nil
}

// TODO:使用gorm造成不能遍历密码，从而验证密钥正确性，需要修改，所以，后期存储密钥直接存储加密后的值
func (s *DbUtils) GetTunnelByMd5Password(p string) (tunnel *models.Tunnel) {
	s.GDb.Model(models.Tunnel{}).Where("Password = ?", p).First(tunnel)
	return tunnel
}

func (s *DbUtils) DelTunnel(id int64) error {
	s.GDb.Model(models.Tunnel{}).Delete(models.Tunnel{Id: id})
	return nil
}

func (s *DbUtils) GetTaskById(id int64) (tunnel *models.Tunnel, e error) {
	if s.GDb.Model(models.Tunnel{}).Where("id = ?").First(tunnel).RowsAffected < 1 {
		e = errors.New(fmt.Sprintf("Tunnel id = %d not found", id))
	}
	return
}

// 离线false,上线true
func (s *DbUtils) GetAllTunnelCountByStatus(status bool, mode string) (count int64) {
	db := s.GDb.Model(models.Tunnel{}).Where("connected = ? ", status)
	if common.IsTunnelMode(mode) {
		db = db.Where("mode = ?", mode)
	}
	db.Count(&count)
	return
}

func (s *DbUtils) GetTunnelListByClientId(valid, clientId int64) ([]models.Tunnel, int) {
	var cli []models.Tunnel
	db := s.GDb.Model(models.Tunnel{}).Where("client_id = ?", clientId)
	if valid == 0 || valid == 1 {
		db = db.Where("valid = ?", valid)
	}
	db.Find(&cli)
	return cli, len(cli)
}
