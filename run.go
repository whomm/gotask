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
	"sync/atomic"
	"os/signal"
	"os"
	"fmt"
	"syscall"
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

	pm := md.PManage{0}
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

	fmt.Println("runhere")

	//创建监听退出chan
    c := make(chan os.Signal)
    //监听指定信号 ctrl+c kill
    signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
    go func() {
        for s := range c {
            switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				atomic.AddInt32(&pm.NowStop,1)
                fmt.Println("退出", s)
                ExitFunc()
            case syscall.SIGUSR1:
                fmt.Println("usr1", s)
            case syscall.SIGUSR2:
                fmt.Println("usr2", s)
            default:
                fmt.Println("other", s)
            }
        }
    }()



}

func ExitFunc()  {
    os.Exit(0)
}
