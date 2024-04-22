package service

import (
	"errors"
	"fmt"
	"xnps/pkg/models"
)

type Tunnel struct {
	*Base
}

func NewTunnel(db *Base) *Tunnel {
	a := &Tunnel{}
	a.Base = db
	return a
}

func (s *Tunnel) GetTunnelListByClientIdWithPage(start, recordsCount int, mode string, clientId int64) ([]models.Tunnel, int) {
	var cli []models.Tunnel
	db := s.Orm(models.Tunnel{}).Where("client_id = ? and valid = 1", clientId)
	//if common.IsTunnelMode(mode) {
	//	db = db.Where("mode = ?", mode)
	//}
	db.Offset(start).Limit(recordsCount).Find(&cli)
	return cli, len(cli)
}
func (s *Tunnel) GetAllTunnelList(status int) ([]models.Tunnel, int) {
	var cli []models.Tunnel
	db := s.Orm(models.Tunnel{})
	if status == 0 || status == 1 {
		db = db.Where("valid = ?", status)
	}
	db.Find(&cli)
	return cli, len(cli)
}
func (s *Tunnel) GetAllTunnelNumById(id int64) uint32 {
	var num int64
	s.Orm(models.Tunnel{}).Where("id = ?", id).Count(&num)
	return uint32(num)
}

func (s *Tunnel) NewTunnel(tunnel *models.Tunnel) (err error) {
	if tunnel.ClientId <= 0 {
		return fmt.Errorf("所属客户端不存在，请重新提交")
	}
	if err != nil {
		return
	}
	s.Orm(models.Tunnel{}).Create(&tunnel)
	return
}
func (s *Tunnel) UpdateTunnel(tunnel *models.Tunnel) error {
	if tunnel.Id <= 0 {
		return fmt.Errorf("tunnel id 填写错误，请检查")
	} else if tunnel.ClientId <= 0 {
		return fmt.Errorf("所属客户端不存在，请重新提交")
	}
	s.Orm(models.Tunnel{}).Where("id = ?", tunnel.Id).Omit().Updates(tunnel)
	return nil
}

// TODO:使用gorm造成不能遍历密码，从而验证密钥正确性，需要修改，所以，后期存储密钥直接存储加密后的值
func (s *Tunnel) GetTunnelByMd5Password(p string) (tunnel *models.Tunnel) {
	s.Orm(models.Tunnel{}).Where("Password = ?", p).First(tunnel)
	return tunnel
}

func (s *Tunnel) DelTunnel(id int64) error {
	s.Orm(models.Tunnel{}).Delete(models.Tunnel{Id: id})
	return nil
}

func (s *Tunnel) GetTaskById(id int64) (tunnel *models.Tunnel, e error) {
	if s.Orm(models.Tunnel{}).Where("id = ?").First(tunnel).RowsAffected < 1 {
		e = errors.New(fmt.Sprintf("Tunnel id = %d not found", id))
	}
	return
}

// 离线false,上线true
func (s *Tunnel) GetAllTunnelCountByStatus(status bool, mode string) (count int64) {
	db := s.Orm(models.Tunnel{}).Where("connected = ? ", status)
	//if common.IsTunnelMode(mode) {
	//	db = db.Where("mode = ?", mode)
	//}
	db.Count(&count)
	return
}

func (s *Tunnel) GetTunnelListByClientId(valid, clientId int64) ([]models.Tunnel, int) {
	var cli []models.Tunnel
	db := s.Orm(models.Tunnel{}).Where("client_id = ?", clientId)
	if valid == 0 || valid == 1 {
		db = db.Where("valid = ?", valid)
	}
	db.Find(&cli)
	return cli, len(cli)
}
