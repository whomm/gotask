package controllers

import (
	"strconv"
	"strings"
	"time"

	hprose "github.com/hprose/hprose-go"

	"encoding/json"

	"../../modle"
)

type MainController struct {
	baseController
}

func (c *MainController) Index() {
	c.Data["adminid"] = c.userid
	c.Data["Username"] = c.username + "/" + c.mygroup.Name
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
func (c *MainController) Test() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "test.tpl"
}

func (this *MainController) Login() {
	if this.IsPost() {
		name := strings.TrimSpace(this.GetString("name"))
		password := strings.TrimSpace(this.GetString("password"))
		var info string
		if len(name) == 0 || len(password) == 0 {
			info = "请填写登录帐号和密码..."
		} else {
			var user modle.User
			user.Name = name
			if user.Read("name") != nil || user.Password != modle.Md5(password) {
				info = "帐号或密码错误..."
			} else {
				/*
					user.Logintimes += 1
					user.Lastloginip = this.getClientIp()
					user.Lastlogintime = this.getTime()
					user.Update()
				*/
				authKey := modle.Md5("whomm|" + user.Password)
				this.Ctx.SetCookie("auth", strconv.FormatInt(user.Id, 10)+"|"+authKey)
				this.Redirect("/", 302)
			}
		}
		this.Data["name"] = name
		this.Data["info"] = info
	}
	this.TplName = "login.tpl"
}

func (this *MainController) Logout() {

	this.Ctx.SetCookie("auth", "")
	this.Ctx.WriteString("<script>top.location.href='" + this.loginurl + "'</script>")
}

func (this *MainController) TaskGroupList() {

	var (
		page     int64
		pagesize int64 = 3
		offset   int64
		list     []*modle.TaskGroup
		task     modle.TaskGroup
		keyword  string
		pager    string

		id int64
	)
	keyword = strings.TrimSpace(this.GetString("keyword"))
	id, _ = this.GetInt64("id")

	pagestr := this.Ctx.Input.Param(":page")
	page, _ = strconv.ParseInt(pagestr, 10, 64)
	if page < 1 {
		page = 1
	}
	offset = (page - 1) * pagesize
	query := task.Query()
	query = query.Filter("ugid__eq", this.mygroup.Id)
	if len(keyword) > 0 {
		query = query.Filter("name__icontains", keyword)
	}
	if id > 0 {
		query = query.Filter("id", id)
	}

	count, _ := query.Count()
	if count > 0 {
		query.OrderBy("-Id").Limit(pagesize, offset).All(&list)
	}

	pager = this.PageList(pagesize, page, count, false, "/taskgrouplist", this.Ctx.Request.URL.RawQuery)
	this.Data["pager"] = pager
	this.Data["pageinfo"] = this.Pageinfo(pagesize, page, count, false, "/taskgrouplist")
	this.Data["pagearr"] = this.PageArr(pagesize, page, count, false, "/taskgrouplist", this.Ctx.Request.URL.RawQuery)
	this.Data["list"] = list
	this.Data["keyword"] = keyword
	this.Data["id"] = id
	this.Data["count"] = count
	this.TplName = "taskgroup.tpl"
}

func (this *MainController) TaskList() {

	var (
		page     int64
		pagesize int64 = 3
		offset   int64
		list     []*modle.Task
		task     modle.Task
		keyword  string
		pager    string

		id int64
	)
	keyword = strings.TrimSpace(this.GetString("keyword"))
	id, _ = this.GetInt64("id")

	pagestr := this.Ctx.Input.Param(":page")
	page, _ = strconv.ParseInt(pagestr, 10, 64)
	if page < 1 {
		page = 1
	}
	offset = (page - 1) * pagesize
	query := task.Query()
	query = query.Filter("ugid__eq", this.mygroup.Id)
	if len(keyword) > 0 {
		query = query.Filter("name__icontains", keyword)
	}
	if id > 0 {
		query = query.Filter("id", id)
	}

	count, _ := query.Count()
	if count > 0 {
		query.OrderBy("-Id").Limit(pagesize, offset).All(&list)
	}

	pager = this.PageList(pagesize, page, count, false, "/tasklist", this.Ctx.Request.URL.RawQuery)
	this.Data["pager"] = pager
	this.Data["pageinfo"] = this.Pageinfo(pagesize, page, count, false, "/tasklist")
	this.Data["pagearr"] = this.PageArr(pagesize, page, count, false, "/tasklist", this.Ctx.Request.URL.RawQuery)
	this.Data["list"] = list
	this.Data["keyword"] = keyword
	this.Data["id"] = id
	this.Data["count"] = count
	this.TplName = "task.tpl"
}

type ServerFunc struct {
	RpcCreateTask func(str string) (int64, error)
}

func (this *MainController) TaskSave() {

	var item modle.Task
	item.Id, _ = this.GetInt64("id")
	item.Uid = this.userid
	item.Tgid, _ = this.GetInt64("tgid")
	item.Ugid = this.mygroup.Id
	item.Name = this.GetString("name")
	item.Crontab = this.GetString("crontab")
	item.Nextrun = 0

	tm, _ := time.Parse("2006-01-02 15:04:05", this.GetString("starttime"))
	item.Starttime = uint(tm.Unix())

	tm, _ = time.Parse("2006-01-02 15:04:05", this.GetString("endtime"))
	item.Endtime = uint(tm.Unix())

	ext := make(map[string]string)
	extval := this.GetStrings("extval[]")
	extkey := this.GetStrings("extkey[]")
	for i, v := range extkey {
		ext[v] = extval[i]
	}
	extb, _ := json.Marshal(ext)
	item.Extra = string(extb)
	item.Invalid, _ = this.GetInt("invalid")
	item.Relay = this.GetString("relay")

	item.Relay = `{"rl":[]}`

	ib, _ := json.Marshal(item)

	//callback finish
	client := hprose.NewClient("http://127.0.0.1:8911/")
	var ro *ServerFunc
	client.UseService(&ro)
	id, err := ro.RpcCreateTask(string(ib))
	if id > 0 && err == nil {
		this.Redirect("/tasklist", 302)
	} else {
		this.showmsg(err.Error())
	}

}

func (this *MainController) TaskUpdate() {

	var tg modle.TaskGroup
	var tglist []*modle.TaskGroup
	tglist = tg.GetAllByUgid(this.mygroup.Id)
	this.Data["tglist"] = tglist
	this.TplName = "taskcreate.tpl"
}

func (this *MainController) TaskInsList() {

	var (
		page     int64
		pagesize int64 = 20
		offset   int64
		list     []*modle.Instance
		taskins  modle.Instance
		task     modle.Task
		taskinfo *modle.Task
		keyword  string
		pager    string

		tid int64
	)
	keyword = strings.TrimSpace(this.GetString("keyword"))
	tid, _ = this.GetInt64("tid")

	pagestr := this.Ctx.Input.Param(":page")
	page, _ = strconv.ParseInt(pagestr, 10, 64)
	if page < 1 {
		page = 1
	}
	offset = (page - 1) * pagesize
	query := taskins.Query()
	if len(keyword) > 0 {
		query = query.Filter("name__icontains", keyword)
	}
	if tid > 0 {
		query = query.Filter("tid", tid)
	}

	count, _ := query.Count()
	if count > 0 {
		query.OrderBy("-Id").Limit(pagesize, offset).All(&list)
	}

	taskinfo = task.GetTaskbyid(tid)

	pager = this.PageList(pagesize, page, count, false, "/taskinslist", this.Ctx.Request.URL.RawQuery)
	this.Data["pager"] = pager
	this.Data["pageinfo"] = this.Pageinfo(pagesize, page, count, false, "/taskinslist")
	this.Data["pagearr"] = this.PageArr(pagesize, page, count, false, "/taskinslist", this.Ctx.Request.URL.RawQuery)
	this.Data["list"] = list
	this.Data["keyword"] = keyword
	this.Data["tid"] = tid
	this.Data["count"] = count
	this.Data["taskinfo"] = taskinfo
	this.TplName = "taskins.tpl"
}

func (this *MainController) UserGroupList() {

	var (
		page     int64
		pagesize int64 = 3
		offset   int64
		list     []*modle.UserGroup
		task     modle.UserGroup
		keyword  string
		pager    string

		id int64
	)
	keyword = strings.TrimSpace(this.GetString("keyword"))
	id, _ = this.GetInt64("id")

	pagestr := this.Ctx.Input.Param(":page")
	page, _ = strconv.ParseInt(pagestr, 10, 64)
	if page < 1 {
		page = 1
	}
	offset = (page - 1) * pagesize
	query := task.Query()
	if len(keyword) > 0 {
		query = query.Filter("name__icontains", keyword)
	}
	if id > 0 {
		query = query.Filter("id", id)
	}

	count, _ := query.Count()
	if count > 0 {
		query.OrderBy("-Id").Limit(pagesize, offset).All(&list)
	}

	pager = this.PageList(pagesize, page, count, false, "/usergrouplist", this.Ctx.Request.URL.RawQuery)
	this.Data["pager"] = pager
	this.Data["pageinfo"] = this.Pageinfo(pagesize, page, count, false, "/usergrouplist")
	this.Data["pagearr"] = this.PageArr(pagesize, page, count, false, "/usergrouplist", this.Ctx.Request.URL.RawQuery)
	this.Data["list"] = list
	this.Data["keyword"] = keyword
	this.Data["id"] = id
	this.Data["count"] = count
	this.TplName = "usergroup.tpl"
}

func (this *MainController) UserList() {

	var (
		page     int64
		pagesize int64 = 20
		offset   int64
		list     []*modle.User
		task     modle.User
		keyword  string
		pager    string

		id int64
	)
	keyword = strings.TrimSpace(this.GetString("keyword"))
	id, _ = this.GetInt64("id")

	pagestr := this.Ctx.Input.Param(":page")
	page, _ = strconv.ParseInt(pagestr, 10, 64)
	if page < 1 {
		page = 1
	}
	offset = (page - 1) * pagesize
	query := task.Query()
	if len(keyword) > 0 {
		query = query.Filter("name__icontains", keyword)
	}
	if id > 0 {
		query = query.Filter("id", id)
	}

	count, _ := query.Count()
	if count > 0 {
		query.OrderBy("-Id").Limit(pagesize, offset).All(&list)
	}

	pager = this.PageList(pagesize, page, count, false, "/userlist", this.Ctx.Request.URL.RawQuery)
	this.Data["pager"] = pager
	this.Data["pageinfo"] = this.Pageinfo(pagesize, page, count, false, "/userlist")
	this.Data["pagearr"] = this.PageArr(pagesize, page, count, false, "/userlist", this.Ctx.Request.URL.RawQuery)
	this.Data["list"] = list
	this.Data["keyword"] = keyword
	this.Data["id"] = id
	this.Data["count"] = count
	this.TplName = "user.tpl"
}
