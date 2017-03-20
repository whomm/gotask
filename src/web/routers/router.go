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
	beego.Router("/logout", &controllers.MainController{}, "*:Logout")

	beego.Router("/tasklist", &controllers.MainController{}, "*:TaskList")
	beego.Router("/tasklist/:page:int/", &controllers.MainController{}, "*:TaskList")
	beego.Router("/taskupdate", &controllers.MainController{}, "*:TaskUpdate")
	beego.Router("/taskinslist", &controllers.MainController{}, "*:TaskInsList")
	beego.Router("/taskinslist/:page:int/", &controllers.MainController{}, "*:TaskInsList")

	beego.Router("/taskgrouplist", &controllers.MainController{}, "*:TaskGroupList")
	beego.Router("/taskgrouplist/:page:int/", &controllers.MainController{}, "*:TaskGroupList")

	//api json
	jns := beego.NewNamespace("/api",
		beego.NSRouter("/tglist", &controllers.ApiController{}, "*:Tglist"),
	)
	beego.AddNamespace(jns)

}
