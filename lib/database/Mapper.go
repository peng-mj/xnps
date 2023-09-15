package database

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"gorm.io/gorm"
	"os"
	"sort"
	"strconv"
	"sync"
	"xnps/lib/common"
	"xnps/lib/crypt"
	"xnps/lib/database/models"
	"xnps/lib/rate"
)

type DbUtils struct {
	GDb *gorm.DB
	//JsonDb *JsonDb
}

func (s *DbUtils) CheckVKey(vKey string) bool {
	return s.GDb.Model(models.Client{}).Where("verify_key = ?", vKey).First(new(models.Client)).RowsAffected > 0
}

func NewClient(vKey string, noStore bool, noDisplay bool) *models.Client {
	return &models.Client{
		VerifyKey:  vKey,
		RemoteAddr: "",
		Name:       "",
		Valid:      true,
		Connected:  false,
		RateLimit:  0,
		Flow:       new(models.Flow),
		Rate:       nil,
		RWMutex:    sync.RWMutex{},
	}
}

// init csv from file
func GetDb() *DbUtils {
	if Db == nil {
		logs.Info("数据库未打开")

		os.Exit(-1)
	}
	return Db
}

func GetMapKeys(m sync.Map, isSort bool, sortKey, order string) (keys []int) {
	if sortKey != "" && isSort {
		return sortClientByKey(m, sortKey, order)
	}
	m.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(int))
		return true
	})
	sort.Ints(keys)
	return
}

func (s *DbUtils) CheckUserName(username string) bool {
	if len(username) > 3 {
		c := int64(0)
		s.GDb.Model(models.SystemConfig{}).Where("web_username = ?", username).Limit(1).Count(&c)
		return c > 0
	}
	return false

}

func (s *DbUtils) CheckTunnelClient(clientId, tunId int64) bool {
	var count int64
	s.GDb.Model(models.Client{}).Where("client_id = ? and id = ? ", clientId, tunId).Count(&count)
	return count > 0
}

func (s *DbUtils) GetClientList(start, length int64, search, sort, order string, clientId int) ([]models.Client, int) {
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

func (s *DbUtils) SetClientStatus(status bool, clientId int64) {
	s.GDb.Model(models.Client{}).Where("id = ?", clientId).Updates(&models.Client{Connected: status})
}

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
func (s *DbUtils) GetAllTunnelNumById(id int64) int {
	var num int64
	s.GDb.Model(models.Tunnel{}).Where("id = ?", id).Count(&num)
	return int(num)
}

func (s *DbUtils) GetIdByVerifyKey(vKey string, addr string) (id int64, err error) {
	var cli models.Client
	res := s.GDb.Model(models.Client{}).Where("verify_key = ?", vKey).RowsAffected
	if res > 0 {
		return cli.Id, nil
	}
	return 0, errors.New("not found")
}

func (s *DbUtils) NewTask(t *models.Tunnel) (err error) {
	if t.Mode == "secret" || t.Mode == "p2p" {
		if s.GDb.Model(models.Tunnel{}).Where("passwd = ?", t.Password).First(new(models.Tunnel)).RowsAffected > 0 {
			err = errors.New(fmt.Sprintf("secret mode keys %s must be unique", t.Password))
		}
	}

	//s.JsonDb.Tasks.Range(func(key, value interface{}) bool {
	//	v := value.(*models.Tunnel)
	//	if (v.Mode == "secret" || v.Mode == "p2p") && v.Password == t.Password {
	//		err = errors.New(fmt.Sprintf("secret mode keys %s must be unique", t.Password))
	//		return false
	//	}
	//	return true
	//})
	if err != nil {
		return
	}
	t.Flow = new(models.Flow)
	s.GDb.Model(models.Tunnel{}).Create(&t)
	return
}

func (s *DbUtils) UpdateTask(t *models.Tunnel) error {
	s.GDb.Model(models.Tunnel{}).Where("id = ?", t.Id).Updates(t)
	return nil
}

func (s *DbUtils) DelTask(id int64) error {
	s.GDb.Model(models.Tunnel{}).Delete(models.Tunnel{Id: id})
	return nil
}

// TODO:使用gorm造成不能遍历密码，从而验证密钥正确性，需要修改，所以，后期存储密钥直接存储加密后的值
func (s *DbUtils) GetTaskByMd5Password(p string) (tunnel *models.Tunnel) {
	//var tunnel models.Tunnel
	s.GDb.Model(models.Tunnel{}).Where("Password = ?", p).First(tunnel)
	return tunnel
	//s.JsonDb.Tasks.Range(func(key, value interface{}) bool {
	//	if crypt.Sha256(value.(*models.Tunnel).Password) == p {
	//		//t = value.(*models.Tunnel)
	//		return false
	//	}
	//	return true
	//})
	//return
}

func (s *DbUtils) GetTaskById(id int64) (tunnel *models.Tunnel, e error) {
	if s.GDb.Model(models.Tunnel{}).Where("id = ?").First(tunnel).RowsAffected < 1 {
		e = errors.New(fmt.Sprintf("Tunnel id = %d not found", id))
	}
	return
}

func (s *DbUtils) DelClient(id int64) error {
	s.GDb.Model(models.Client{}).Delete(&models.Client{Id: id})
	return nil
}

// 这个地方需要重构，客户端ID自动生成
func (s *DbUtils) NewClient(c *models.Client) error {
	var isNotSet bool
	if c.HttpUser != "" && !s.VerifyUserName(c.HttpUser, c.Id) {
		return errors.New("web login username duplicate, please reset")
	}
reset:
	if c.VerifyKey == "" || isNotSet {
		isNotSet = true
		c.VerifyKey = crypt.GenerateRandomVKey()
	}
	if c.RateLimit == 0 {
		c.Rate = rate.NewRate(int64(2 << 23))
	} else if c.Rate == nil {
		c.Rate = rate.NewRate(int64(c.RateLimit * 1024))
	}
	c.Rate.Start()
	if !s.VerifyVkey(c.VerifyKey, c.Id) {
		if isNotSet {
			goto reset
		}
		return errors.New("vkey duplicate, please reset")
	}
	//去掉，让系统自动生成
	//if c.Id == 0 {
	//	c.Id = s.JsonDb.GetClientId()
	//}
	if c.Flow == nil {
		c.Flow = new(models.Flow)
	}
	s.GDb.Model(models.Client{}).Create(&c)
	return nil
}

func (s *DbUtils) VerifyVkey(vkey string, id int64) (res bool) {
	return s.GDb.Model(models.Client{}).Where("verify_key = ? AND id = ?", vkey, id).First(new(models.Client)).RowsAffected > 0
}

func (s *DbUtils) VerifyUserName(username string, id int64) (res bool) {
	return s.GDb.Model(models.Client{}).Where("web_user = ? AND id = ?", username, id).First(new(models.Client)).RowsAffected > 0
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

// 检查是否启用
func (s *DbUtils) IsPubClient(id int64) bool {
	return s.GDb.Model(models.Client{}).Where("id = ? and valid = 1", id).First(new(models.Client)).RowsAffected > 0
}

func (s *DbUtils) GetClientById(id int64) (c *models.Client, err error) {
	if s.GDb.Model(models.Client{}).Where("id = ? and valid = 1", id).First(c).RowsAffected < 1 {
		err = errors.New("未找到客户端")
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
func (s *DbUtils) HasTunnel(clientId int64, t *models.Tunnel) bool {
	var num int64
	s.GDb.Model(models.Tunnel{}).Where("client_id = ?", clientId).Where(t).Count(&num)
	return num > 0
	//return s.GDb.Model(models.Tunnel{}).Where("client_id = ?", clientId).Where(t).First(new(models.Tunnel)).RowsAffected > 0
}

func (s *DbUtils) GetClientTunnelNumByClientId(id int64) int {
	var num int64
	s.GDb.Model(models.Tunnel{}).Where("client_id = ?", id).Count(&num)
	return int(num)
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
func (s *DbUtils) EditSysConfig(config *models.SystemConfig) {

	s.GDb.Model(models.SystemConfig{}).Updates(config)
}
func (s *DbUtils) GetSystemConfig() (sys models.SystemConfig, err error) {
	if s.GDb.Model(models.SystemConfig{}).First(&sys).RowsAffected < 1 {
		err = errors.New("have no sys config")
	}
	return
}
