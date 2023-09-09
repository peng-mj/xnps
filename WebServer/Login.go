package WebServer

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"io"
	"net/http"
	"xnps/WebServer/WebObj"
	"xnps/lib/crypt"
)

func (w *WebServer) keyAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		// 获取请求中的密钥
		//key := c.Request().Header.Get("token")
		//user := c.Request().Header.Get("username")
		//Utils.Log("user:", user)
		//Utils.Log("token:", key)
		//Utils.Log("URL:", c.Request().URL)
		//Utils.Log("header中的token:" + key)
		//token, err := Database.KvWebGetToken(user)

		// 验证密钥是否正确
		//if err != nil || key != token {
		//Utils.Log("token 验证失败")
		//re := WebObj.ReData{Data: ReCode.ErTokenLapsed, Status: "ERROR"}
		//request, _ := json.Marshal(re)
		//return c.HTML(http.StatusAccepted, string(request))
		//}
		//Utils.Log("密钥验证通过")
		// 密钥验证通过，继续传递请求给下一个处理程序
		return next(c)
	}
}

// Login 登录信息
func (w *WebServer) Login(c echo.Context) (err error) {
	re := WebObj.Request{}
	body, err := io.ReadAll(c.Request().Body)
	login := new(WebObj.Login)
	if err == nil {
		json.Unmarshal(body, &re)
		if data, ok := re.Data.(map[string]interface{}); ok {
			err = mapstructure.Decode(data, login)
			passwd, _ := w.salt.Get(login.Username)
			if login.Username != "root" && login.Password != passwd {
				err = errors.New("username or password error")
			}
		}
	}
	if err != nil {
		re.MsgType = "ERROR"
		re.Data = err.Error()
	} else {
		re.MsgType = "OK"
		re.Data = "登录成功"
	}
	request, _ := json.Marshal(re)
	return c.HTML(http.StatusOK, string(request))
}

// DoLogin 用于登录前准备，当用户登录时，将Salt和密码的sha256进行组合后再sha256编码，再传回服务器，避免中间者拦截密码的sha256
func (w *WebServer) DoLogin(c echo.Context) (err error) {
	re := WebObj.Request{}
	body, err := io.ReadAll(c.Request().Body)
	//doLogin := new(WebObj.DoLogin)
	var Salt string
	if err == nil {
		json.Unmarshal(body, &re)
		if userName, ok := re.Data.(string); ok {
			//TODO:这里需要验证用户是否在数据库中
			Salt = crypt.GenerateRandomVKey()
			if err == nil {
				err = w.salt.Add(userName, Salt)
			}
		}
	}

	if err != nil {
		re.MsgType = "ERROR"
		re.Data = err.Error()
	} else {
		re.MsgType = "OK"
		re.Data = Salt
	}

	request, _ := json.Marshal(re)
	return c.HTML(http.StatusOK, string(request))
}
