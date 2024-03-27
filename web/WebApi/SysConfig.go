package WebApi

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"xnps/database/Mapper"
	"xnps/database/models"
	"xnps/web/WebObj"
)

/**********        USER          *********/

func GetSysConfig(c echo.Context) (err error) {
	var conf models.SystemConfig
	conf, err = Mapper.GetDb().GetSystemConfig()
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
			Mapper.GetDb().UpdateSysConfig(conf)
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
