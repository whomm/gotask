package routers

import (
	"../controllers"
	"github.com/astaxie/beego"
)

func init() {

	//pages
	beego.Router("/", &controllers.MainController{}, "*:Index")
	beego.Router("/test", &controllers.MainController{}, "*:Test")
	beego.Router("/login", &controllers.MainController{}, "*:Login")

	//api json
	jns := beego.NewNamespace("/api",
		beego.NSRouter("/tglist", &controllers.ApiController{}, "*:Tglist"),
	)
	beego.AddNamespace(jns)

}
