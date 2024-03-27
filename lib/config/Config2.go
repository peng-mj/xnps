package config

import (
	"github.com/go-ini/ini"
	"log/slog"
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
