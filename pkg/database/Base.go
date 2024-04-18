package database

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"xnps/pkg/common"
	"xnps/pkg/database/models"
	"xnps/web/service"
)

var Db *service.DbUtils

func InitDatabase(dbFile string) *service.DbUtils {
	var db = service.DbUtils{}
	var err error
	db.GDb, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		fmt.Println("打开数据库失败")
		os.Exit(-1)
	} else {
		err := db.GDb.AutoMigrate(
			models.Client{}, models.Tunnel{}, models.Group{}, models.Blacklist{})
		if err != nil {
			logs.Info("创建数据表失败", err)
		}
	}
	Db = &db
	return &db
}

func NewJsonDb(runPath string) *JsonDb {
	return &JsonDb{
		RunPath:        runPath,
		TaskFilePath:   filepath.Join(runPath, "conf", "tasks.json"),
		ClientFilePath: filepath.Join(runPath, "conf", "clients.json"),
	}
}

type JsonDb struct {
	Tasks            sync.Map
	Clients          sync.Map
	RunPath          string
	ClientIncreaseId int32  //client increased id
	TaskIncreaseId   int32  //task increased id
	TaskFilePath     string //task file path
	ClientFilePath   string //client file path
}

//func (s *JsonDb) LoadTaskFromJsonFile() {
//	loadSyncMapFromFile(s.TaskFilePath, func(v string) {
//		var err error
//		post := new(models.Tunnel)
//		if json.Unmarshal([]byte(v), &post) != nil {
//			return
//		}
//		if post.Client, err = s.GetClientById(post.Client.Id); err != nil {
//			return
//		}
//		s.Tasks.Store(post.Id, post)
//		if post.Id > int(s.TaskIncreaseId) {
//			s.TaskIncreaseId = int32(post.Id)
//		}
//	})
//}

//func (s *JsonDb) LoadClientFromJsonFile() {
//	loadSyncMapFromFile(s.ClientFilePath, func(v string) {
//		post := new(models.Client)
//		if json.Unmarshal([]byte(v), &post) != nil {
//			return
//		}
//		if post.RateLimit > 0 {
//			post.Rate = rate.NewRate(int64(post.RateLimit * 1024))
//		} else {
//			post.Rate = rate.NewRate(int64(2 << 23))
//		}
//		post.Rate.Start()
//		post.NowConn = 0
//		s.Clients.Store(post.Id, post)
//		if post.Id > s.ClientIncreaseId) {
//			s.ClientIncreaseId = int32(post.Id)
//		}
//	})
//}

//func (s *JsonDb) GetClientById(id int64) (c *models.Client, err error) {
//	var cli models.Client
//
//	if v, ok := s.Clients.Load(id); ok {
//		c = v.(*models.Client)
//		return
//	}
//	err = errors.New("未找到客户端")
//	return
//}

//var hostLock sync.Mutex

//
//func (s *JsonDb) StoreHostToJsonFile() {
//	hostLock.Lock()
//	storeSyncMapToFile(s.Hosts, s.HostFilePath)
//	hostLock.Unlock()
//}

//var taskLock sync.Mutex

//func (s *JsonDb) StoreTasksToJsonFile() {
//	taskLock.Lock()
//	storeSyncMapToFile(s.Tasks, s.TaskFilePath)
//	taskLock.Unlock()
//}

//var clientLock sync.Mutex

//func (s *JsonDb) StoreClientsToJsonFile() {
////	clientLock.Lock()
////	storeSyncMapToFile(s.Clients, s.ClientFilePath)
////	clientLock.Unlock()
////}

func (s *JsonDb) GetClientId() int32 {
	return atomic.AddInt32(&s.ClientIncreaseId, 1)
}

func (s *JsonDb) GetTaskId() int32 {
	return atomic.AddInt32(&s.TaskIncreaseId, 1)
}

func loadSyncMapFromFile(filePath string, f func(value string)) {
	b, err := common.ReadAllFromFile(filePath)
	if err != nil {
		panic(err)
	}
	for _, v := range strings.Split(string(b), "\n"+common.CONN_DATA_SEQ) {
		f(v)
	}
}

//
//func storeSyncMapToFile(m sync.Map, filePath string) {
//	file, err := os.Create(filePath + ".tmp")
//	// first create a temporary file to store
//	if err != nil {
//		panic(err)
//	}
//	m.Range(func(key, value interface{}) bool {
//		var b []byte
//		var err error
//		switch value.(type) {
//		case *models.Tunnel:
//			obj := value.(*models.Tunnel)
//			if obj.NoStore {
//				return true
//			}
//			b, err = json.Marshal(obj)
//		//case *models.Host:
//		//	obj := value.(*models.Host)
//		//	if obj.NoStore {
//		//		return true
//		//	}
//		//	b, err = json.Marshal(obj)
//		case *models.Client:
//			obj := value.(*models.Client)
//			if obj.Valid {
//				return true
//			}
//			b, err = json.Marshal(obj)
//		default:
//			return true
//		}
//		if err != nil {
//			return true
//		}
//		_, err = file.Write(b)
//		if err != nil {
//			panic(err)
//		}
//		_, err = file.Write([]byte("\n" + common.CONN_DATA_SEQ))
//		if err != nil {
//			panic(err)
//		}
//		return true
//	})
//	_ = file.Sync()
//	_ = file.Close()
//	// must close file first, then rename it
//	err = os.Rename(filePath+".tmp", filePath)
//	if err != nil {
//		logs.Error(err, "store to file err, data will lost")
//	}
//	// replace the file, maybe provides atomic operation
//}
