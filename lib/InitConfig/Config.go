package InitConfig

import (
	"github.com/go-ini/ini"
	"golang.org/x/exp/slog"
	"os"
	"strconv"
	"strings"
)

type IniConf struct {
	cfg  *ini.File
	path string
}

var Ini *IniConf

func NewIniConfig(path string) *IniConf {
	if f, err := ini.Load(path); err != nil {
		slog.Error("config file don't exist, please check it. path=" + path)
		os.Exit(-1)
	} else {
		Ini = &IniConf{
			cfg:  f,
			path: path,
		}
	}
	return Ini
}

func CheckIni() bool {
	return Ini.cfg != nil
}

func GetIni() *IniConf {
	if Ini.cfg == nil {
		slog.Error("config file didn't open or didn't have it")
		os.Exit(-1)
	}
	return Ini
}

func (i *IniConf) GetValue(sec string, key string) string {

	iniValue := i.cfg.Section(sec).Key(key).Value()
	return iniValue
}

func (i *IniConf) GetString(sec string, key string) string {
	return i.cfg.Section(sec).Key(key).Value()
}
func (i *IniConf) GetInt(sec string, key string) int {
	v := i.GetString(sec, key)
	out, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return out
}
func (i *IniConf) GetFloat(sec string, key string) float64 {
	v := i.GetString(sec, key)
	out, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0
	}
	return out
}

// 当目标是1、true、ok、yes、y时为真
func (i *IniConf) GetBool(sec string, key string) bool {
	v := strings.ToUpper(i.GetString(sec, key))
	return v == "1" || v == "TRUE" || v == "OK" || v == "YES" || v == "Y"
}

func (i *IniConf) SetString(sec string, key string, value string) {
	i.cfg.Section(sec).Key(key).SetValue(value)
	err := i.cfg.SaveTo(i.path)
	if err != nil {
		return
	}
	//如果flash需要刷新，那么这里需要处理，避免断电丢失
	//Hardware.RunSync()
}

func (i *IniConf) SetInt(sec string, key string, value int) {
	s := strconv.Itoa(value)
	i.SetString(sec, key, s)
}
func (i *IniConf) SetFloat(sec string, key string, value float64) {
	s := strconv.FormatFloat(value, 'f', -1, 64)
	i.SetString(sec, key, s)
}
func (i *IniConf) SetBool(sec string, key string, value bool) {
	s := strconv.FormatBool(value)
	i.SetString(sec, key, s)
}
