package controllers

import (
	"../../modle"
)

type ApiController struct {
	baseController
}

func init() {

}

type OutJson struct {
	Code int64
	Data interface{}
	Msg  string
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

	j := OutJson{}
	j.Code = 0
	j.Data = out
	j.Msg = ""

	this.Data["json"] = j
	this.ServeJSON()

}
