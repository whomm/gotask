package main

import (
	_ "./routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.AppPath = "../"
	beego.Run()
}
