package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/astaxie/beego/logs"
	hprose "github.com/hprose/hprose-go"
)

type myWorker struct{}

type clientStub struct {
	RpcKilledTask   func(int64) error
	RpcFailTask     func(int64) error
	RpcScuccessTask func(int64) error
}

//回调服务器完成任务
func dothejob(taskinstanceid int64) {

	time.Sleep(10)
	//callback finish
	client := hprose.NewClient("http://127.0.0.1:8911/")
	var ro *clientStub
	client.UseService(&ro)
	fmt.Println(ro.RpcScuccessTask(taskinstanceid))

}

//回调服务器杀死任务成功
func killthejob(taskinstanceid int64) {
	time.Sleep(10)
	client := hprose.NewClient("http://127.0.0.1:8911/")
	var ro *clientStub
	client.UseService(&ro)
	fmt.Println(ro.RpcKilledTask(taskinstanceid))
}

func (myWorker) Run(tasktime int64, taskid int64, taskinstanceid int64, taskinfo string) bool {
	logs.Debug(taskinfo)
	log.Println(tasktime, taskid, taskinstanceid)
	go dothejob(taskinstanceid)
	return true
}

func (myWorker) Kill(taskinstanceid int64) bool {
	logs.GetLogger("WORK").Println("kill : %v", taskinstanceid)
	go killthejob(taskinstanceid)
	return true

}

func main() {
	logs.SetLogger(logs.AdapterFile, `{"filename":"worker.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
	service := hprose.NewHttpService()
	service.AddMethods(myWorker{})
	http.ListenAndServe(":8912", service)
}
