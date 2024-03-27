package web

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/juju/ratelimit"
	"net/http"
	"strings"
	"sync"
	"time"
	"xnps/database/Mapper"
	"xnps/lib/cache"
	"xnps/lib/crypt"
	"xnps/web/WebApi"
	"xnps/web/WebObj"
)

type WebServer struct {
	//e        *echo.Echo
	g        *gin.Engine
	salt     *cache.Cache
	JWTToken *jwt.Token
	secret   *cache.Cache
	//tokenMan TokenManager
}

func (w *WebServer) Start(url string) {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Done()
	w.g = gin.Default()
	w.g.Use(gin.Logger(), gin.Recovery())
	w.salt = cache.New(20)
	w.secret = cache.New(20) //默认最大密钥存储的文件数量
	w.JWTToken = jwt.New(jwt.SigningMethodHS256)
	//w.e = echo.New()

	//w.e.HTTPErrorHandler = w.ErrorHandler
	//装载静态文件
	//TODO:后续使用go:embed打包
	w.g.Static("/", "web/static/")
	w.g.POST("/doLogin", w.DoLogin)
	w.g.POST("/login", w.Login)
	w.ApiRouterRules("/api")
	//w.e.Logger.Fatal(w.e.Start(url))
}

func (w *WebServer) jwtMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		// 检查Token是否存在
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			c.Abort() //提前结束
			return
		}
		// 获取Token
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证算法和密钥
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if uuid, ok := claims["uuid"].(string); ok {
						if secret, ok := w.secret.GetString(uuid); ok {
							return []byte(secret), nil
						}
					}
				}
			}
			return nil, errors.New(http.StatusText(http.StatusUnauthorized))
		})
		// Token验证通过，将Token信息存储在Context中，以便后续处理函数使用，这个用作目录管理
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && err == nil {
			//TODO：这个地方可有可无->Set
			c.Set("user", claims)
			return
		}
		c.Next()
		//c.JSON(http.StatusUnauthorized, WebApi.Replay(err, nil))
	}
}

func (w *WebServer) generateToken(username, uid, passwd string, timeoutHour int) string {
	claims := w.JWTToken.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["uid"] = uid
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(timeoutHour)).Unix() // 设置Token的有效期
	tokenString, _ := w.JWTToken.SignedString([]byte(passwd))
	return tokenString
}
func (w *WebServer) RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.String(http.StatusForbidden, "rate limit...")
			c.Abort()
			return
		}
		c.Next()
	}
}

// Login 登录信息
func (w *WebServer) Login(c *gin.Context) {
	var login WebObj.Login
	var token string
	var err error

	if err = c.ShouldBindJSON(&login); err == nil {
		if Salt, ok := w.salt.GetString(login.Username); ok {
			//slog.Info("salt=", Salt)
			var passwd string
			if passwd, err = Mapper.GetDb().GetPasswdByUser(login.Username); err == nil {
				//验证加密后的密码是否正确
				if passwd = crypt.Sha256(Salt + "." + passwd); login.Password != passwd {
					err = errors.New("username or password error")
				} else {
					salt := crypt.GetUlid()
					token = w.generateToken(login.Username, salt, passwd, 1)
					w.secret.Add(salt, passwd)
					w.salt.Remove(login.Username)
				}
			}
		} else {
			err = errors.New(http.StatusText(http.StatusUnauthorized))
		}
	}
	c.JSON(http.StatusOK, WebApi.Replay(err, token))
}

// DoLogin 用于登录前准备，当用户登录时，将Salt和密码的sha256进行组合后再sha256编码，再传回服务器，避免中间者拦截密码的sha256
func (w *WebServer) DoLogin(c *gin.Context) {
	var doLogin WebObj.DoLogin
	var Salt string
	var err error
	if err = c.ShouldBindJSON(&doLogin); err == nil {
		if Mapper.GetDb().CheckUserName(doLogin.Username) {
			Salt = crypt.GenerateRandomVKey()
			w.salt.Add(doLogin.Username, Salt)
		}
	}
	c.JSON(http.StatusOK, WebApi.Replay(err, Salt))
}
