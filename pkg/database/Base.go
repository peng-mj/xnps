package database

import (
	"errors"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	SQLITE = 1
	MYSQL  = 2
	SQL    = 3
)

type Driver struct {
	db         gorm.Dialector
	orm        *gorm.DB
	dbType     int
	dbPath     string
	tableModel []interface{}
}

func New() *Driver {
	return &Driver{}
}
func (c *Driver) Sqlite(name string) *Driver {
	c.dbPath = name
	c.dbType = SQLITE
	return c
}

// Mysql Not yet enabled
// url=username:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
func (c *Driver) Mysql(url string) *Driver {
	c.dbPath = url
	c.dbType = MYSQL
	return c
}
func (c *Driver) AddModels(dst []interface{}) *Driver {
	c.tableModel = dst
	return c
}

func (c *Driver) Init() (driver *Driver, err error) {
	if len(c.tableModel) == 0 {
		return nil, errors.New("please add the table structure")
	}
	switch c.dbType {
	case SQLITE:
		c.db = sqlite.Open(c.dbPath)
	case MYSQL:
		c.db = mysql.Open(c.dbPath)
	default:
		return nil, errors.New("database config not init")
	}
	c.orm, err = gorm.Open(c.db, &gorm.Config{})
	if err == nil {
		err = c.orm.AutoMigrate(c.tableModel...)
		if err != nil {
			err = errors.New("failed to create or update the table structure")
		}
	}
	return c, err
}
func (c *Driver) Orm(t interface{}) *gorm.DB {
	return c.orm.Model(t)
}
