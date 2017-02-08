package main

import (
	"net/http"

	md "./modle"
	_ "github.com/go-sql-driver/mysql"
	hprose "github.com/hprose/hprose-go"
)

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

	http.ListenAndServe(":8911", service)

}
