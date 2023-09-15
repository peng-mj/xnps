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
	return c.HTML(http.StatusOK, string(reStr))
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
	return c.HTML(http.StatusOK, string(reStr))
}

// TODO：配置文件还是单独存入一个配置文件中，通过监测配置文件是否存在来启动初始化启动器
// 添加新的系统配置，如果某些配置为空，那么使用默认配置，将生成的默认配置返回
func AddSysConfig(c echo.Context) (err error) {
	re := WebObj.Request{}
	body, err := io.ReadAll(c.Request().Body)
	conf := new(models.SystemConfig)
	if err == nil {
		err = json.Unmarshal(body, &conf)
		if err == nil {
			conf, err = database.GetDb().AddSysConfig(conf)
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
	//return c.HTML(http.StatusOK, string(reStr))
	return c.HTML(http.StatusOK, string(reStr))
}

//func CheckSysInit(c echo.Context) (err error) {
//
//}

/**********        GROUP          *********/

func GetAllGroup(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func GetGroupByCondition(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func AddGroup(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func DelGroup(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func EditGroup(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}

/**********        CLIENT          *********/

func GetAllClient(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func GetClientByCondition(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func AddClient(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func DelClient(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func EditClient(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}

/**********        TUNNEL          *********/

func GetAllTunnel(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func GetTunnelByCondition(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func AddTunnel(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func DelTunnel(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func EditTunnel(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}

/**********        FIREWALL          *********/

func GetAllFirewall(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func GetGroupByFirewall(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func AddFirewall(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func DelFirewall(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func EditFirewall(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}

/**********        SYSTEM          *********/

func GetSystemStatus(c echo.Context) (err error) {
	return c.HTML(http.StatusOK, "")
}
func GetConnectInfo(c echo.Context) (err error) {
	return c.HTML(http.StatusOK, "")
}

/**********        BLOCK_LIST          *********/

func GetAllBlockList(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func GetBlockListByCondition(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func AddFBlockList(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func DelBlockList(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
func EditBlockList(c echo.Context) (err error) {

	return c.HTML(http.StatusOK, "")
}
