package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	hprose "github.com/hprose/hprose-go"
)

type myWorker struct{}

func dothejob(taskinstanceid int64) {

	time.Sleep(10)
	//callback finish

	esp, err := http.Get("http://192.168.0.253:8034/v1api/insstatsfish?id=" + string(taskinstanceid))
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

}
func killthejob(taskinstanceid int64) {
	time.Sleep(10)
	//callback

	esp, err := http.Get("http://192.168.0.253:8034/v1api/insstatskilled?id=" + string(taskinstanceid))
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// ...
}

func (myWorker) Run(tasktime int64, taskid int64, taskinstanceid int64, taskinfo string) bool {
	log.Println(taskinfo)
	log.Println(tasktime, taskid, taskinstanceid)
	go dothejob()
	return true
}

func (myWorker) Kill(taskinstanceid int64) bool {
	log.Println("kill : %v", taskinstanceid)
	go killthejob()
	return true

}

func main() {
	service := hprose.NewHttpService()
	service.AddMethods(myWorker{})
	http.ListenAndServe(":8912", service)
}
