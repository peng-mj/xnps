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
	orm        *gorm.DB
	dbType     int
	dbPath     string
	tableModel []interface{}
}

func New() *Driver {
	return &Driver{}
}
func (c *Driver) NewSqlite(name string) {
	c.dbPath = name
	c.dbType = SQLITE
}

// NewMysql Not yet enabled
// url=username:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
func (c *Driver) NewMysql(url string) *Driver {
	c.dbPath = url
	c.dbType = MYSQL
	return c
}
func (c *Driver) AddTable(dst ...interface{}) {
	c.tableModel = dst
}

func (c *Driver) Init() (driver *Driver, err error) {
	if len(c.tableModel) == 0 {
		return nil, errors.New("please add the table structure")
	}
	switch c.dbType {
	case SQLITE:
		c.orm, err = gorm.Open(sqlite.Open(c.dbPath), &gorm.Config{})
	case MYSQL:
		c.orm, err = gorm.Open(mysql.Open(c.dbPath), &gorm.Config{})
	default:
		return nil, errors.New("database config not init")
	}
	if err == nil {
		err = c.orm.AutoMigrate(c.tableModel)
		if err != nil {
			err = errors.New("failed to create or update the table structure")
		}
	}
	return c, err
}
func (c *Driver) Orm(t interface{}) *gorm.DB {
	return c.orm.Model(t)
}
