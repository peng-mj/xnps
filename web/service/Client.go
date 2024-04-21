package service

import (
	"errors"
	"strconv"
	"xnps/pkg/database"
	"xnps/pkg/models"
	"xnps/pkg/rate"
)

type Client struct {
	Base
}

func NewClient(db *database.Driver) *Client {
	c := &Client{}
	c.Service(db)
	return c
}

func (c *Client) SetConnectedStatus(status bool, clientId int64) {
	c.Orm(models.Client{}).Where("id = ?", clientId).Updates(&models.Client{Connected: status})
}

func (c *Client) CreatClient(vKey string) *models.Client {
	return &models.Client{
		AccessKey:  vKey,
		RemoteAddr: "",
		Name:       "",
		Valid:      true,
		Connected:  false,
		RateLimit:  0,
	}
}

func (c *Client) GetAllClientList(start, length int64, search, sort, order string, clientId int) ([]models.Client, int) {
	var cli []models.Client
	c.GDb.Model(models.Client{}).Where("valid = 1").Find(&cli)
	return cli, len(cli)
}

// GetAllClientCount if status <0 this func would count whether the client connected, the status equal to 1 for connected and o for disconnected clients
// the valid is like to status
func (c *Client) GetAllClientCount(status int, valid int) (count int64) {
	db := c.GDb.Model(models.Tunnel{})
	if status == 0 || status == 1 {
		db = db.Where("connected = ?", status)
	}
	if valid == 0 || valid == 1 {
		db = db.Where("valid = ?", valid)
	}
	db.Count(&count)
	return
}

func (c *Client) GetClientCountByMode(mode string) (count int64) {
	c.GDb.Model(models.Client{}).Where("mode = ? ", mode).Count(&count)
	return
}
func (c *Client) UpdateClientById(client *models.Client, id int64) (count int64) {
	c.GDb.Model(models.Client{}).Where("id = ?", id).Updates(client)
	return
}
func (c *Client) GetClientByAccessUser(AccessId, AccessKey string) (client models.Client, err error) {
	if c.GDb.Model(models.Client{}).Where("access_id = ? and access_key = ?", AccessId, AccessKey).First(&client).RowsAffected == 0 {
		err = errors.New("client not found")
	}
	return
}

func (c *Client) GetClientIdByVkey(vkey string) (id int64, err error) {
	var cli models.Client
	err = errors.New("未找到客户端")
	if c.GDb.Model(models.Client{}).Where("verify_key = ?", vkey).First(&cli).RowsAffected < 1 {
	}
	return
}

func (c *Client) GetClientTunnelNumByClientId(id int64) int {
	var num int64
	c.GDb.Model(models.Tunnel{}).Where("client_id = ?", id).Count(&num)
	return int(num)
}

func (c *Client) HasTunnel(clientId int64, t *models.Tunnel) bool {
	var num int64
	c.GDb.Model(models.Tunnel{}).Where("client_id = ?", clientId).Where(t).Count(&num)
	return num > 0
	//return c.GDb.Model(models.Tunnel{}).Where("client_id = ?", clientId).Where(t).First(new(models.Tunnel)).RowsAffected > 0
}
func (c *Client) GetClientById(id int64) (client *models.Client, err error) {
	if c.GDb.Model(models.Client{}).Where("id = ? and valid = 1", id).First(client).RowsAffected < 1 {
		err = errors.New("未找到客户端")
	}
	return
}

// 检查是否启用
func (c *Client) CheckClientValid(id int64) bool {
	var count int64
	return c.GDb.Model(models.Client{}).Where("id = ? and valid = 1", id).Limit(1).Count(&count).RowsAffected > 0
}
func (c *Client) UpdateClient(t *models.Client) error {
	//c.JsonDb.Clients.Store(t.Id, t)

	res := c.GDb.Model(models.Client{}).Where("id = ?", t.Id).Updates(t).RowsAffected
	if t.RateLimit == 0 {
		t.Rate = rate.NewRate(int64(2 << 23))
		t.Rate.Start()
	}
	if res < 1 {
		return errors.New("have no client or the client have no change where id =  " + strconv.FormatInt(t.Id, 10))
	}
	return nil
}

func (c *Client) VerifyUserName(username string, id int64) (res bool) {
	return c.GDb.Model(models.Client{}).Where("web_user = ? AND id = ?", username, id).First(new(models.Client)).RowsAffected > 0
}

func (c *Client) CreateNewClient(client *models.Client) error {
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
	if !c.VerifyClientVkey(client.AccessKey, client.Id) {
		if isNotSet {
			goto reset
		}
		return errors.New("vkey duplicate, please reset")
	}

	c.GDb.Model(models.Client{}).Create(client)
	return nil
}
func (s *DbUtils) DelClient(id int64) error {
	s.GDb.Model(models.Client{}).Delete(&models.Client{Id: id})
	return nil
}
func (s *DbUtils) VerifyClientVkey(vkey string, id int64) (res bool) {
	return s.GDb.Model(models.Client{}).Where("verify_key = ? AND id = ?", vkey, id).First(new(models.Client)).RowsAffected > 0
}
