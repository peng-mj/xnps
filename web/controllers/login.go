package controllers

import (
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/utils/captcha"
	"math/rand"
	"net"
	"sync"
	"time"
	"xnps/lib/database/models"

	"github.com/astaxie/beego"
	"xnps/lib/common"
	"xnps/lib/database"
	"xnps/server"
)

type LoginController struct {
	beego.Controller
}

var ipRecord sync.Map
var cpt *captcha.Captcha

type record struct {
	hasLoginFailTimes int
	lastLoginTime     time.Time
}

func init() {
	// use beego cache system store the captcha data
	store := cache.NewMemoryCache()
	cpt = captcha.NewWithFilter("/captcha/", store)
}

func (self *LoginController) Index() {
	// Try login implicitly, will succeed if it's configured as no-auth(empty username&password).
	webBaseUrl := beego.AppConfig.String("web_base_url")
	if self.doLogin("", "", false) {
		self.Redirect(webBaseUrl+"/index/index", 302)
	}
	self.Data["web_base_url"] = webBaseUrl
	self.Data["register_allow"], _ = beego.AppConfig.Bool("allow_user_register")
	self.Data["captcha_open"], _ = beego.AppConfig.Bool("open_captcha")
	self.TplName = "login/index.html"
}

func (self *LoginController) Verify() {
	username := self.GetString("username")
	password := self.GetString("password")
	captchaOpen, _ := beego.AppConfig.Bool("open_captcha")
	if captchaOpen {
		if !cpt.VerifyReq(self.Ctx.Request) {
			self.Data["json"] = map[string]interface{}{"status": 0, "msg": "the verification code is wrong, please get it again and try again"}
			self.ServeJSON()
		}
	}
	if self.doLogin(username, password, true) {
		self.Data["json"] = map[string]interface{}{"status": 1, "msg": "login success"}
	} else {
		self.Data["json"] = map[string]interface{}{"status": 0, "msg": "username or password incorrect"}
	}
	self.ServeJSON()
}

// TODO:登录逻辑需要重构
// 原本可能是为了实现多用户，但是，没有完全实现，留下了后门，但凡获得了一个客户端的权限
func (self *LoginController) doLogin(username, password string, explicit bool) bool {
	clearIpRecord()
	ip, _, _ := net.SplitHostPort(self.Ctx.Request.RemoteAddr)
	if v, ok := ipRecord.Load(ip); ok {
		vv := v.(*record)
		if (time.Now().Unix() - vv.lastLoginTime.Unix()) >= 60 {
			vv.hasLoginFailTimes = 0
		}
		if vv.hasLoginFailTimes >= 10 {
			return false
		}
	}
	var auth bool
	if password == beego.AppConfig.String("web_password") && username == beego.AppConfig.String("web_username") {
		self.SetSession("isAdmin", true)
		self.DelSession("clientId")
		self.DelSession("username")
		auth = true
		server.Bridge.Register.Store(common.GetIpByAddr(self.Ctx.Input.IP()), time.Now().Add(time.Hour*time.Duration(2)))
	}
	b, err := beego.AppConfig.Bool("allow_user_login")
	//TODO:这个地方需要理解一下
	//允许客户用户登录
	//这个地方去掉，不允许多用户登录
	if err == nil && b && !auth {
		//database.GetDb().JsonDb.Clients.Range(func(key, value interface{}) bool {
		//	v := value.(*models.Client)
		//	if !v.Status || v.Valid {
		//		return true
		//	}
		//	if v.HttpUser == "" && v.HttpPasswd == "" {
		//		//为什么是 VerifyKey！=password
		//		if username != "user" || v.VerifyKey != password {
		//			return true
		//		} else {
		//			auth = true
		//		}
		//	}
		//	if !auth && v.HttpUser == password && v.HttpPasswd == username {
		//		auth = true
		//	}
		//	if auth {
		//		self.SetSession("isAdmin", false)
		//		self.SetSession("clientId", v.Id)
		//		self.SetSession("username", v.HttpUser)
		//		return false
		//	}
		//	return true
		//})
	}
	if auth {
		self.SetSession("auth", true)
		ipRecord.Delete(ip)
		return true

	}
	if v, load := ipRecord.LoadOrStore(ip, &record{hasLoginFailTimes: 1, lastLoginTime: time.Now()}); load && explicit {
		vv := v.(*record)
		vv.lastLoginTime = time.Now()
		vv.hasLoginFailTimes += 1
		ipRecord.Store(ip, vv)
	}
	return false
}
func (self *LoginController) Register() {
	if self.Ctx.Request.Method == "GET" {
		self.Data["web_base_url"] = beego.AppConfig.String("web_base_url")
		self.TplName = "login/register.html"
	} else {
		if b, err := beego.AppConfig.Bool("allow_user_register"); err != nil || !b {
			self.Data["json"] = map[string]interface{}{"status": 0, "msg": "register is not allow"}
			self.ServeJSON()
			return
		}
		if self.GetString("username") == "" || self.GetString("password") == "" || self.GetString("username") == beego.AppConfig.String("web_username") {
			self.Data["json"] = map[string]interface{}{"status": 0, "msg": "please check your input"}
			self.ServeJSON()
			return
		}
		t := &models.Client{
			//Id:         database.GetDb().JsonDb.GetClientId(),
			Valid: true,
			//Cnf:        &models.Config{},
			HttpUser:   self.GetString("username"),
			HttpPasswd: self.GetString("password"),
			Flow:       &models.Flow{},
		}
		if err := database.GetDb().NewClient(t); err != nil {
			self.Data["json"] = map[string]interface{}{"status": 0, "msg": err.Error()}
		} else {
			self.Data["json"] = map[string]interface{}{"status": 1, "msg": "register success"}
		}
		self.ServeJSON()
	}
}

func (self *LoginController) Out() {
	self.SetSession("auth", false)
	self.Redirect(beego.AppConfig.String("web_base_url")+"/login/index", 302)
}

func clearIpRecord() {
	x := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100)
	if x == 1 {
		ipRecord.Range(func(key, value interface{}) bool {
			v := value.(*record)
			if time.Now().Unix()-v.lastLoginTime.Unix() >= 60 {
				ipRecord.Delete(key)
			}
			return true
		})
	}
}
