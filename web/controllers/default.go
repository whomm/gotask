package controllers

import (
	"strconv"
	"strings"
	"time"

	hprose "github.com/hprose/hprose-go"

	"encoding/json"

	"github.com/whomm/gotask/modle"
)

const TASKCTLURL string = "http://127.0.0.1:8911/"

type MainController struct {
	baseController
}

//任务主页面
func (c *MainController) Index() {
	c.Data["adminid"] = c.userid
	c.Data["Username"] = c.username + "/" + c.mygroup.Name
	c.Data["Website"] = "http://github.com/whomm"
	c.Data["Email"] = "lacing@126.com"
	c.TplName = "index.tpl"
}

//
func (c *MainController) Test() {
	c.Data["Website"] = "http://github.com/whomm"
	c.Data["Email"] = "lacing@126.com"
	c.TplName = "test.tpl"
}

//用户登录页面
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

//登出页
func (this *MainController) Logout() {

	this.Ctx.SetCookie("auth", "")
	this.Ctx.WriteString("<script>top.location.href='" + this.loginurl + "'</script>")
}

//任务组列表页
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

//任务列表页
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

//任务创建和更新的rpc接口
type ServerFunc struct {
	RpcCreateTask func(str string) (int64, error)
	RpcUpdateTask func(str string, fields string) (bool, error)
}

//更新或创建任务操作
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
		if len(strings.Trim(v, "\n ")) == 0 {
			continue
		}
		ext[v] = extval[i]
	}
	extb, _ := json.Marshal(ext)
	item.Extra = string(extb)
	item.Invalid, _ = this.GetInt("invalid")
	item.Relay = this.GetString("relay")

	item.Relay = `{"rl":[]}`

	ib, _ := json.Marshal(item)

	//callback finish
	client := hprose.NewClient(TASKCTLURL)
	var ro *ServerFunc
	client.UseService(&ro)

	if item.Id < 1 {
		id, err := ro.RpcCreateTask(string(ib))
		if id > 0 && err == nil {
			this.Redirect("/tasklist", 302)
		} else {
			this.showmsg(err.Error())
		}
	} else {
		id, err := ro.RpcUpdateTask(string(ib), "")
		if id && err == nil {
			this.Redirect("/tasklist", 302)
		} else {
			this.showmsg(err.Error())
		}
	}

}

//任务创建或更新页面
func (this *MainController) TaskUpdate() {

	var tg modle.TaskGroup
	var tglist []*modle.TaskGroup
	var taskinfo *modle.Task
	tglist = tg.GetAllByUgid(this.mygroup.Id)
	this.Data["tglist"] = tglist
	taskid, ok := this.GetInt64("id")
	if ok == nil {
		var tk modle.Task
		taskinfo = tk.GetTaskbyid(taskid)

		this.Data["taskinfo"] = taskinfo
		ext := make(map[string]string)
		ok = json.Unmarshal([]byte(taskinfo.Extra), &ext)
		type ExtX struct {
			Key string
			Val string
		}
		var extarray []ExtX
		if ok == nil {
			for k, v := range ext {
				extarray = append(extarray, ExtX{Key: k, Val: v})
			}
		}

		this.Data["taskinfoext"] = extarray
		this.TplName = "taskupdate.tpl"
	} else {
		this.TplName = "taskcreate.tpl"
	}

}

//任务运行记录或实例列表
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

//用户组列表
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

//用户列表页
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
