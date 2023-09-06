package database

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"gorm.io/gorm"
	"os"
	"sort"
	"sync"
	"xnps/lib/crypt"
	"xnps/lib/database/models"
	"xnps/lib/rate"
)

type DbUtils struct {
	GDb    *gorm.DB
	JsonDb *JsonDb
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

func (s *DbUtils) GetClientList(start, length int, search, sort, order string, clientId int) ([]models.Client2, int) {
	var cli []models.Client2
	s.GDb.Model(models.Client2{}).Where("valid=1").Find(&cli)
	return cli, len(cli)

}

func (s *DbUtils) GetIdByVerifyKey(vKey string, addr string) (id int, err error) {
	var cli models.Client2
	res := s.GDb.Model(models.Client2{}).Where("verify_key = ?", vKey).RowsAffected
	if res > 0 {
		return int(cli.Id), nil
	}
	return 0, errors.New("not found")
}

func (s *DbUtils) NewTask(t *models.Tunnel) (err error) {
	s.JsonDb.Tasks.Range(func(key, value interface{}) bool {
		v := value.(*models.Tunnel)
		if (v.Mode == "secret" || v.Mode == "p2p") && v.Password == t.Password {
			err = errors.New(fmt.Sprintf("secret mode keys %s must be unique", t.Password))
			return false
		}
		return true
	})
	if err != nil {
		return
	}
	t.Flow = new(models.Flow)
	s.JsonDb.Tasks.Store(t.Id, t)
	s.JsonDb.StoreTasksToJsonFile()
	return
}

func (s *DbUtils) UpdateTask(t *models.Tunnel) error {
	s.JsonDb.Tasks.Store(t.Id, t)
	s.JsonDb.StoreTasksToJsonFile()
	return nil
}

func (s *DbUtils) DelTask(id int) error {
	s.JsonDb.Tasks.Delete(id)
	s.JsonDb.StoreTasksToJsonFile()
	return nil
}

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

func (s *DbUtils) GetTask(id int) (t *models.Tunnel, err error) {
	if v, ok := s.JsonDb.Tasks.Load(id); ok {
		t = v.(*models.Tunnel)
		return
	}
	err = errors.New("not found")
	return
}

func (s *DbUtils) DelClient(id int) error {
	s.JsonDb.Clients.Delete(id)
	s.JsonDb.StoreClientsToJsonFile()
	return nil
}

func (s *DbUtils) NewClient(c *models.Client) error {
	var isNotSet bool
	if c.WebUserName != "" && !s.VerifyUserName(c.WebUserName, c.Id) {
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
		return errors.New("Vkey duplicate, please reset")
	}
	if c.Id == 0 {
		c.Id = int(s.JsonDb.GetClientId())
	}
	if c.Flow == nil {
		c.Flow = new(models.Flow)
	}
	s.JsonDb.Clients.Store(c.Id, c)
	s.JsonDb.StoreClientsToJsonFile()
	return nil
}

func (s *DbUtils) VerifyVkey(vkey string, id int) (res bool) {
	res = true
	s.JsonDb.Clients.Range(func(key, value interface{}) bool {
		v := value.(*models.Client)
		if v.VerifyKey == vkey && v.Id != id {
			res = false
			return false
		}
		return true
	})
	return res
}

func (s *DbUtils) VerifyUserName(username string, id int) (res bool) {
	res = true
	s.JsonDb.Clients.Range(func(key, value interface{}) bool {
		v := value.(*models.Client)
		if v.WebUserName == username && v.Id != id {
			res = false
			return false
		}
		return true
	})
	return res
}

func (s *DbUtils) UpdateClient(t *models.Client) error {
	s.JsonDb.Clients.Store(t.Id, t)
	if t.RateLimit == 0 {
		t.Rate = rate.NewRate(int64(2 << 23))
		t.Rate.Start()
	}
	return nil
}

// 检查是否启用
func (s *DbUtils) IsPubClient(id int) bool {
	client, err := s.GetClient(id)
	if err == nil {
		return client.NoDisplay
	}
	return false
}

func (s *DbUtils) GetClient(id int) (c *models.Client, err error) {
	if v, ok := s.JsonDb.Clients.Load(id); ok {
		c = v.(*models.Client)
		return
	}
	err = errors.New("未找到客户端")
	return
}

func (s *DbUtils) GetClientIdByVkey(vkey string) (id int, err error) {
	var exist bool
	s.JsonDb.Clients.Range(func(key, value interface{}) bool {
		v := value.(*models.Client)
		if crypt.Md5(v.VerifyKey) == vkey {
			exist = true
			id = v.Id
			return false
		}
		return true
	})
	if exist {
		return
	}
	err = errors.New("未找到客户端")
	return
}
