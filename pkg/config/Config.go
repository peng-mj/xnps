package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/go-ini/ini"
	"golang.org/x/exp/slog"
	"io/ioutil"
	"strconv"
	"strings"
)

type Driver struct {
	cfg         *ini.File
	iniFilePath string
}

func (s *Driver) GetValue(sec string, key string) string {
	iniValue := s.cfg.Section(sec).Key(key).Value()
	return iniValue
}
func (s *Driver) GetInt(sec string, key string) int {
	iniValue := s.cfg.Section(sec).Key(key).Value()
	out, err := strconv.ParseInt(iniValue, 10, 32)
	if err != nil {
		return 0
	}
	return int(out)
}
func (s *Driver) GetInt64(sec string, key string, def int64) int64 {
	str := s.GetValue(sec, key)
	if i, err := strconv.ParseInt(str, 10, 64); err == nil {
		return i
	}
	return def
}
func (s *Driver) GetBool(sec string, key string, def bool) bool {
	v := s.GetValue(sec, key)
	if v == "" {
		return def
	}
	v = strings.ToLower(v)
	return v == "1" || v == "true" || v == "t" || v == "ok" || v == "yes" || v == "y"
}

func (s *Driver) SetValue(sec string, key string, value string) {
	s.cfg.Section(sec).Key(key).SetValue(value)
	err := s.cfg.SaveTo(s.iniFilePath)
	if err != nil {
		slog.Error("配置文件保存失败")
		return
	}
}
func (s *Driver) SetBool(sec string, key string, value bool) {
	v := ""
	if value {
		v = "true"
	} else {
		v = "false"
	}
	s.cfg.Section(sec).Key(key).SetValue(v)
	err := s.cfg.SaveTo(s.iniFilePath)
	if err != nil {
		slog.Error("配置文件保存失败")
		return
	}
}
func (s *Driver) SetInt(sec string, key string, value int64) {
	s.cfg.Section(sec).Key(key).SetValue(strconv.FormatInt(value, 10))
	err := s.cfg.SaveTo(s.iniFilePath)
	if err != nil {
		slog.Error("配置文件保存失败")
		return
	}
}
func NewSysConfig(iniFile string) (cfg *Driver, err error) {
	cfg = new(Driver)
	if cfg.cfg, err = ini.Load(iniFile); err != nil {
		return
	}
	cfg.iniFilePath = iniFile
	return
}

// 假设你有一个配置结构体
type Config struct {
	Name string `toml:"name"`
	Age  int    `toml:"age"`
	// 可以添加更多字段和嵌套结构
}

// Query 查询配置
func Query(filePath string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// AddOrModify 新增或修改配置
func AddOrModify(filePath string, config *Config) error {
	data, err := toml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, data, 0644)
}

// Delete 删除配置中的字段（示例为删除Name字段）
func Delete(filePath string, fieldName string) error {
	config, err := Query(filePath)
	if err != nil {
		return err
	}

	// 根据字段名删除字段（这里仅处理顶级字段，对于嵌套字段需要更复杂的逻辑）
	switch fieldName {
	case "name":
		config.Name = ""
	case "age":
		config.Age = 0
	// 添加更多字段的删除逻辑...
	default:
		return fmt.Errorf("unknown field: %s", fieldName)
	}

	return AddOrModify(filePath, config)
}
