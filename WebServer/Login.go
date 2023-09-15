package WebServer

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"xnps/WebServer/WebObj"
	"xnps/lib/crypt"
	"xnps/lib/database"
)

// func (w *WebServer) keyAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
//
//		return func(c echo.Context) error {
//
//			// 获取请求中的密钥
//			//key := c.Request().Header.Get("token")
//			//user := c.Request().Header.Get("username")
//			//Utils.Log("user:", user)
//			//Utils.Log("token:", key)
//			//Utils.Log("URL:", c.Request().URL)
//			//Utils.Log("header中的token:" + key)
//			//token, err := Database.KvWebGetToken(user)
//
//			// 验证密钥是否正确
//			//if err != nil || key != token {
//			//Utils.Log("token 验证失败")
//			//re := WebObj.ReData{Data: ReCode.ErTokenLapsed, Status: "ERROR"}
//			//request, _ := json.Marshal(re)
//			//return c.HTML(http.StatusAccepted, string(request))
//			//}
//			//Utils.Log("密钥验证通过")
//			// 密钥验证通过，继续传递请求给下一个处理程序
//			return next(c)
//		}
//	}
func (w *WebServer) jwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")

		// 检查Token是否存在
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Missing JWT Token",
			})
		}

		// 解析Token
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		log.Println("token =" + tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证算法和密钥
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.ErrUnauthorized
			}
			return []byte("your_secret_key"), nil
		})

		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Invalid JWT Token",
			})
		}

		// Token验证通过，将Token信息存储在Context中，以便后续处理函数使用
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user", claims)
			return next(c)
		}

		return echo.ErrUnauthorized
	}
}

func (w *WebServer) keyAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		//jwtMap := jwt.New()
		//jwtMap

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

func (w *WebServer) generateToken(username string) string {
	claims := w.JWTToken.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // 设置Token的有效期
	tokenString, _ := w.JWTToken.SignedString([]byte("your_secret_key"))
	return tokenString
}

// Login 登录信息
func (w *WebServer) Login(c echo.Context) (err error) {
	re := WebObj.Request{}
	body, err := io.ReadAll(c.Request().Body)
	login := new(WebObj.Login)
	var token string
	if err == nil {
		json.Unmarshal(body, &re)
		if data, ok := re.Data.(map[string]interface{}); ok {
			err = mapstructure.Decode(data, login)
			Salt, err := w.salt.Get(login.Username)
			if err == nil {
				var passwd string
				passwd, err = database.GetDb().GetPasswdByUser(login.Username)
				if err == nil {
					passwd = crypt.Sha256(Salt + passwd) //验证加密后的密码是否正确
					if login.Password != passwd {
						err = errors.New("username or password error")
					} else {
						token = w.generateToken(login.Username)
					}
				}

			}

		}
	}
	if err != nil {
		re.MsgType = "ERROR"
		re.Data = err.Error()
	} else {
		re.MsgType = "OK"
		re.Data = token
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
		if username, ok := re.Data.(string); ok && database.GetDb().CheckUserName(username) {
			Salt = crypt.GenerateRandomVKey()
			if err == nil {
				err = w.salt.Add(username, Salt)
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
