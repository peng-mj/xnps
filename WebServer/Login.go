package WebServer

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
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
					if uuid, ok := claims["uuid"].(string); ok {
						if pa, ok := w.secret.GetString(uuid); ok {
							return []byte(pa), nil
						} else {
							return nil, echo.ErrUnauthorized
						}
					}
				}
			}
			return nil, echo.ErrUnauthorized
		})
		if err != nil {
			return c.String(http.StatusUnauthorized, WebApi.ReDara(err, nil))
		}
		// Token验证通过，将Token信息存储在Context中，以便后续处理函数使用，这个用作目录管理
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user", claims)
			return next(c)
		}

		return echo.ErrUnauthorized
	}
}

func (w *WebServer) generateToken(username, uuid, passwd string, timeoutHour int) string {
	claims := w.JWTToken.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["uuid"] = uuid
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(timeoutHour)).Unix() // 设置Token的有效期
	tokenString, _ := w.JWTToken.SignedString([]byte(passwd))
	return tokenString
}

// Login 登录信息
func (w *WebServer) Login(c echo.Context) (err error) {
	login := new(WebObj.Login)
	var token string
	body, err := io.ReadAll(c.Request().Body)
	if err == nil {
		if err = json.Unmarshal(body, login); err == nil {
			if Salt, ok := w.salt.GetString(login.Username); ok {
				slog.Info("salt=", Salt)
				var passwd string
				if passwd, err = database.GetDb().GetPasswdByUser(login.Username); err == nil {
					//验证加密后的密码是否正确
					if passwd = crypt.Sha256(Salt + "@" + passwd); login.Password != passwd {
						err = errors.New("username or password error")
					} else {
						uuid := crypt.GetUlid()
						token = w.generateToken(login.Username, uuid, passwd, 1)
						w.secret.Add(uuid, passwd)
						w.salt.Remove(login.Username)
					}
				}
			} else {
				err = echo.ErrUnauthorized
			}
		}
	}
	return c.HTML(http.StatusOK, WebApi.ReDara(err, token))
}

// DoLogin 用于登录前准备，当用户登录时，将Salt和密码的sha256进行组合后再sha256编码，再传回服务器，避免中间者拦截密码的sha256
func (w *WebServer) DoLogin(c echo.Context) (err error) {
	body, err := io.ReadAll(c.Request().Body)
	doLogin := new(WebObj.DoLogin)
	var Salt string
	if err == nil {
		if err = json.Unmarshal(body, &doLogin); err == nil {
			if database.GetDb().CheckUserName(doLogin.Username) {
				Salt = crypt.GenerateRandomVKey()
				w.salt.Add(doLogin.Username, Salt)
			}
		}
	}
	return c.String(http.StatusOK, WebApi.ReDara(err, Salt))
}
