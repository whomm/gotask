package main

import (
	"net/http"
	"time"

	md "github.com/whomm/gotask/modle"

	_ "github.com/whomm/gotask/web/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	hprose "github.com/hprose/hprose-go"
)

func int64todate(in int64) (out string) {
	return time.Unix(int64(in), 0).Format("2006-01-02 15:04:05")

}
func uinttodate(in uint) (out string) {
	return time.Unix(int64(in), 0).Format("2006-01-02 15:04:05")

}

func init() {
	logs.SetLogger(logs.AdapterFile, `{"filename":"project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
}

func main() {

	pm := md.PManage{}
	go pm.GeneraterConsume()
	go pm.RunerConsume()

	service := hprose.NewHttpService()
	service.AddFunction("RpcCreateTask", pm.RpcCreateTask)
	service.AddFunction("RpcUpdateTask", pm.RpcUpdateTask)
	service.AddFunction("RpcRemoveTask", pm.RemoveTask)
	service.AddFunction("RpcKillInstance", pm.RpcKillInstance)
	service.AddFunction("RpcCreateTaskInstance", pm.RpcCreateTaskInstance)
	service.AddFunction("RpcScuccessTask", pm.RpcScuccessTask)
	service.AddFunction("RpcFailTask", pm.RpcFailTask)
	service.AddFunction("RpcKilledTask", pm.RpcKilledTask)

	go http.ListenAndServe(":8911", service)

	beego.AddFuncMap("int64todate", int64todate)
	beego.AddFuncMap("uinttodate", uinttodate)
	beego.Run()

}
