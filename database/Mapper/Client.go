package Mapper

import (
	"errors"
	"strconv"
	"sync"
	"xnps/database/models"
	"xnps/lib/crypt"
	"xnps/lib/rate"
)

func (s *DbUtils) SetClientStatus(status bool, clientId int64) {
	s.GDb.Model(models.Client{}).Where("id = ?", clientId).Updates(&models.Client{Connected: status})
}

func NewClient(vKey string) *models.Client {
	return &models.Client{
		AccessKey:  vKey,
		RemoteAddr: "",
		Name:       "",
		Valid:      true,
		Connected:  false,
		RateLimit:  0,
		Rate:       nil,
		RWMutex:    sync.RWMutex{},
	}
}

func (s *DbUtils) GetAllClientList(start, length int64, search, sort, order string, clientId int) ([]models.Client, int) {
	var cli []models.Client
	s.GDb.Model(models.Client{}).Where("valid = 1").Find(&cli)
	return cli, len(cli)
}

// GetAllClientCount if status <0 this func would count whether the client connected, the status equal to 1 for connected and o for disconnected clients
// the valid is like to status
func (s *DbUtils) GetAllClientCount(status int, valid int) (count int64) {
	db := s.GDb.Model(models.Tunnel{})
	if status == 0 || status == 1 {
		db = db.Where("connected = ?", status)
	}
	if valid == 0 || valid == 1 {
		db = db.Where("valid = ?", valid)
	}
	db.Count(&count)
	return
}

func (s *DbUtils) GetClientCountByMode(mode string) (count int64) {
	s.GDb.Model(models.Client{}).Where("mode = ? ", mode).Count(&count)
	return
}
func (s *DbUtils) UpdateClientById(client *models.Client, id int64) (count int64) {
	s.GDb.Model(models.Client{}).Where("id = ?", id).Updates(client)
	return
}
func (s *DbUtils) GetClientByAccessUser(AccessId, AccessKey string) (client models.Client, err error) {
	if s.GDb.Model(models.Client{}).Where("access_id = ? and access_key = ?", AccessId, AccessKey).First(&client).RowsAffected == 0 {
		err = errors.New("client not found")
	}
	return
}

func (s *DbUtils) GetClientIdByVkey(vkey string) (id int64, err error) {
	var cli models.Client
	err = errors.New("未找到客户端")
	if s.GDb.Model(models.Client{}).Where("verify_key = ?", vkey).First(&cli).RowsAffected < 1 {
	}
	return
}

func (s *DbUtils) GetClientTunnelNumByClientId(id int64) int {
	var num int64
	s.GDb.Model(models.Tunnel{}).Where("client_id = ?", id).Count(&num)
	return int(num)
}

func (s *DbUtils) HasTunnel(clientId int64, t *models.Tunnel) bool {
	var num int64
	s.GDb.Model(models.Tunnel{}).Where("client_id = ?", clientId).Where(t).Count(&num)
	return num > 0
	//return s.GDb.Model(models.Tunnel{}).Where("client_id = ?", clientId).Where(t).First(new(models.Tunnel)).RowsAffected > 0
}
func (s *DbUtils) GetClientById(id int64) (c *models.Client, err error) {
	if s.GDb.Model(models.Client{}).Where("id = ? and valid = 1", id).First(c).RowsAffected < 1 {
		err = errors.New("未找到客户端")
	}
	return
}

// 检查是否启用
func (s *DbUtils) CheckClientValid(id int64) bool {
	var count int64
	return s.GDb.Model(models.Client{}).Where("id = ? and valid = 1", id).Limit(1).Count(&count).RowsAffected > 0
}
func (s *DbUtils) UpdateClient(t *models.Client) error {
	//s.JsonDb.Clients.Store(t.Id, t)

	res := s.GDb.Model(models.Client{}).Where("id = ?", t.Id).Updates(t).RowsAffected
	if t.RateLimit == 0 {
		t.Rate = rate.NewRate(int64(2 << 23))
		t.Rate.Start()
	}
	if res < 1 {
		return errors.New("have no client or the client have no change where id =  " + strconv.FormatInt(t.Id, 10))
	}
	return nil
}

// TODO:不允许客户端登录
func (s *DbUtils) VerifyUserName(username string, id int64) (res bool) {
	return s.GDb.Model(models.Client{}).Where("web_user = ? AND id = ?", username, id).First(new(models.Client)).RowsAffected > 0
}

// 这个地方需要重构，客户端ID自动生成
func (s *DbUtils) CreateNewClient(client *models.Client) error {
	var isNotSet bool
reset:
	if client.AccessKey == "" || isNotSet {
		isNotSet = true
		client.AccessKey = crypt.GenerateRandomVKey()
	}
	if client.RateLimit == 0 {
		client.Rate = rate.NewRate(int64(2 << 23))
	} else if client.Rate == nil {
		client.Rate = rate.NewRate(int64(client.RateLimit * 1024))
	}
	client.Rate.Start()
	if !s.VerifyClientVkey(client.AccessKey, client.Id) {
		if isNotSet {
			goto reset
		}
		return errors.New("vkey duplicate, please reset")
	}

	s.GDb.Model(models.Client{}).Create(client)
	return nil
}
func (s *DbUtils) DelClient(id int64) error {
	s.GDb.Model(models.Client{}).Delete(&models.Client{Id: id})
	return nil
}
func (s *DbUtils) VerifyClientVkey(vkey string, id int64) (res bool) {
	return s.GDb.Model(models.Client{}).Where("verify_key = ? AND id = ?", vkey, id).First(new(models.Client)).RowsAffected > 0
}
