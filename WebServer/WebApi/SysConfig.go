package WebApi

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"xnps/WebServer/WebObj"
	"xnps/lib/database"
	"xnps/lib/database/models"
)

/**********        USER          *********/

func GetSysConfig(c echo.Context) (err error) {
	var conf models.SystemConfig
	conf, err = database.GetDb().GetSystemConfig()
	if err != nil {
		return err
	}
	re := WebObj.Request{}
	if err != nil {
		re.MsgType = "ERROR"
		re.Data = conf
	} else {
		re.MsgType = "OK"
		re.Data = err.Error()
	}
	reStr, _ := json.Marshal(re)
	return c.String(http.StatusOK, string(reStr))
}
func EditSYsConfig(c echo.Context) (err error) {
	re := WebObj.Request{}
	body, err := io.ReadAll(c.Request().Body)
	conf := new(models.SystemConfig)
	if err == nil {
		err = json.Unmarshal(body, &conf)
		if err == nil {
			database.GetDb().EditSysConfig(conf)
		}
	}
	if err != nil {
		re.MsgType = "ERROR"
		re.Data = err.Error()
	} else {
		re.MsgType = "OK"
		re.Data = "update system config succeed"
	}
	reStr, _ := json.Marshal(re)
	return c.String(http.StatusOK, string(reStr))
}

// TODO：配置文件还是单独存入一个配置文件中，通过监测配置文件是否存在来启动初始化启动器
// 添加新的系统配置，如果某些配置为空，那么使用默认配置，将生成的默认配置返回
func AddSysConfig(c echo.Context) (err error) {
	body, err := io.ReadAll(c.Request().Body)
	conf := new(models.SystemConfig)
	if err == nil {
		err = json.Unmarshal(body, &conf)
		if err == nil {
			//TODO:这里判断username的合法性
			conf, err = database.GetDb().AddSysConfig(conf)
		}
	}

	return c.String(http.StatusOK, ReDara(err, "update system config succeed"))
}
