package routers

import (
	"github.com/astaxie/beego"
	"xnps/web/controllers"
)

func Init() {
	webBaseUrl := beego.AppConfig.String("web_base_url")
	if len(webBaseUrl) > 0 {
		ns := beego.NewNamespace(webBaseUrl,
			beego.NSRouter("/", &controllers.IndexController{}, "*:Index"),
			beego.NSAutoRouter(&controllers.IndexController{}),
			beego.NSAutoRouter(&controllers.LoginController{}),
			beego.NSAutoRouter(&controllers.ClientController{}),
			beego.NSAutoRouter(&controllers.AuthController{}),
		)
		beego.AddNamespace(ns)
	} else {
		beego.Router("/", &controllers.IndexController{}, "*:Index")
		beego.AutoRouter(&controllers.IndexController{})
		beego.AutoRouter(&controllers.LoginController{})
		beego.AutoRouter(&controllers.ClientController{})
		beego.AutoRouter(&controllers.AuthController{})
	}
}
