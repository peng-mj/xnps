package config

import (
	"fmt"
	"golang.org/x/exp/slog"
	"testing"
	"time"
	"xnps/pkg/crypt"
	"xnps/pkg/sysTool"
)

func Test(t *testing.T) {
	base := "conf/server.toml"
	dr := New(base)
	conf, err := dr.Load()
	if err != nil {
		return
	}
	fmt.Println(conf)
	conf.Remark = "XNps config file"
	conf.WebPort = 8090
	conf.InitTime = time.Now().Unix()
	conf.AppKeys = crypt.RandStr().AddNum().AddLetter().GenerateList(24, 5)
	err = dr.Update(conf)
	if err != nil {
		return
	}

}
func Test2(t *testing.T) {
	base := "E:/Magic/xnps/conf/server2.toml"
	if sysTool.DirExisted("E:/Magic/xnps/conf") {
		sysTool.CreateFolder("E:/Magic/xnps/conf")
	}
	err := CreateNewInitFile(base)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	dr := New(base)
	conf, err := dr.Load()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	conf.WebPort = 8090
	conf.InitTime = time.Now().Unix()
	conf.AppKeys = crypt.RandStr().AddNum().AddLetter().GenerateList(24, 10)
	dr.Update(conf)

}
