package main

import (
	"fmt"

	hprose "github.com/hprose/hprose-go"
)

type clientStub struct {
	RpcCreateTask func(string) (int64, error)
}

func main() {
	client := hprose.NewClient("http://127.0.0.1:8911/")
	var ro *clientStub
	client.UseService(&ro)
	taskjson := `
{
    "uid": 1,
    "ugid": 1,
    "tgid": 1,
    "name": "test1",
    "crontab": "5 * * * * *",
    "pengdingtime": 0,
    "starttime": 0,
    "endtime": 1798736461,
    "extra": "{\"workrpc\":\"http://127.0.0.1:8912/\"}",
    "invalid": 0,
    "relay": "{\"rl\":[]}"
}
    `
	id, err := ro.RpcCreateTask(taskjson)
	fmt.Println(id)
	fmt.Println(err)
}
