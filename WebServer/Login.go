package WebServer

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"io"
	"net/http"
	"strings"
	"time"
	"xnps/WebServer/WebApi"
	"xnps/WebServer/WebObj"
	"xnps/lib/crypt"
	"xnps/lib/database"
)

func (w *WebServer) jwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		// 检查Token是否存在
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Missing JWT Token",
			})
		}
		// 获取Token
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证算法和密钥
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if username, ok := claims["username"].(string); ok {
						se, err := w.secret.Get(username)
						if err != nil {
							return nil, err
						} else {
							return []byte(se), nil
						}
					}
				}
			}
			return nil, echo.ErrUnauthorized
		})
		if err != nil {
			return c.String(http.StatusUnauthorized, WebApi.ReDara(err, nil))
		}
		// Token验证通过，将Token信息存储在Context中，以便后续处理函数使用
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user", claims)
			return next(c)
		}

		return echo.ErrUnauthorized
	}
}

// 目前没有什么作用
func (w *WebServer) keyAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		// 密钥验证通过，继续传递请求给下一个处理程序
		return next(c)
	}
}

func (w *WebServer) generateToken(username string) string {
	claims := w.JWTToken.Claims.(jwt.MapClaims)
	claims["username"] = username + crypt.GenerateRandomVKey()
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
		err = json.Unmarshal(body, &login)
		if err != nil {
			return err
		}
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
