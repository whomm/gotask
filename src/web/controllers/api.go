package controllers

import (
	"../../modle"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type ApiController struct {
	beego.Controller
}

func init() {

	//设置最大空闲连接
	//设置最大数据库连接 (go >= 1.2)
	maxIdle := 30
	maxConn := 30
	// set default database
	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("mysqlurls"), maxIdle, maxConn)
	orm.Debug = true
}

type TaskGropPage struct {
	List []*modle.TaskGroup
	Ps   int64
	Pl   int64
}

func (this *ApiController) Tglist() {

	var (
		taskgroup *modle.TaskGroup = new(modle.TaskGroup)
		out       *TaskGropPage    = new(TaskGropPage)
	)
	this.Ctx.Output.Header("Cache-Control", "public")

	out.List = taskgroup.GetList(0, 10)
	out.Ps = 0
	out.Pl = 10

	this.Data["json"] = out
	this.ServeJSON()

}
