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
	"xnps/lib/crypt"
	"xnps/lib/database/models"
	"xnps/lib/rate"
)

type DbUtils struct {
	GDb    *gorm.DB
	JsonDb *JsonDb
}

func (s *DbUtils) CheckVKey(vKey string) bool {
	return s.GDb.Model(models.Client{}).Where("verify_key = ?", vKey).First(new(models.Client)).RowsAffected > 0
}

func NewClient(vKey string, noStore bool, noDisplay bool) *models.Client {
	return &models.Client{
		VerifyKey: vKey,
		Addr:      "",
		Remark:    "",
		Valid:     true,
		Connected: false,
		RateLimit: 0,
		Flow:      new(models.Flow),
		Rate:      nil,
		RWMutex:   sync.RWMutex{},
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

func (s *DbUtils) GetClientList(start, length int, search, sort, order string, clientId int) ([]models.Client, int) {
	var cli []models.Client
	s.GDb.Model(models.Client{}).Where("valid = 1").Find(&cli)
	return cli, len(cli)

}

func (s *DbUtils) GetIdByVerifyKey(vKey string, addr string) (id int, err error) {
	var cli models.Client
	res := s.GDb.Model(models.Client{}).Where("verify_key = ?", vKey).RowsAffected
	if res > 0 {
		return int(cli.Id), nil
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

// TODO:后期修改为sha验证
// md5 password
func (s *DbUtils) GetTaskByMd5Password(p string) (t *models.Tunnel) {
	s.JsonDb.Tasks.Range(func(key, value interface{}) bool {
		if crypt.Md5(value.(*models.Tunnel).Password) == p {
			t = value.(*models.Tunnel)
			return false
		}
		return true
	})
	return
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

func (s *DbUtils) NewClient(c *models.Client) error {
	var isNotSet bool
	if c.WebUser != "" && !s.VerifyUserName(c.WebUser, c.Id) {
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
func (s *DbUtils) IsPubClient(id int) bool {
	return s.GDb.Model(models.Client{}).Where("id = ? and valid = 1", id).First(new(models.Client)).RowsAffected > 0
}

func (s *DbUtils) GetClient(id int) (c *models.Client, err error) {
	if s.GDb.Model(models.Client{}).Where("id = ? and valid = 1", id).First(c).RowsAffected < 1 {
		err = errors.New("未找到客户端")
	}
	return
}

func (s *DbUtils) GetClientIdByVkey(vkey string) (id int64, err error) {
	var cli models.Client
	if s.GDb.Model(models.Client{}).Where("verify_key = ?", vkey).First(&cli).RowsAffected < 1 {
		err = errors.New("未找到客户端")
	}
	return

}
