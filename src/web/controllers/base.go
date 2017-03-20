package controllers

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"../../modle"

	"log"

	"github.com/astaxie/beego"
)

type baseController struct {
	beego.Controller
	userid         int64
	username       string
	nickname       string
	mygroup        modle.UserGroup
	controllerName string
	actionName     string
	loginurl       string
}

func (this *baseController) Prepare() {
	controllerName, actionName := this.GetControllerAndAction()
	this.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
	this.actionName = strings.ToLower(actionName)
	this.auth()
}

///用户权限验证
func (this *baseController) auth() {
	this.loginurl = beego.AppConfig.String("loginurl")
	if this.actionName != "login" && this.actionName != "logout" {
		ck := strings.Split(this.Ctx.GetCookie("auth"), "|")
		if len(ck) == 2 {
			idstr, password := ck[0], ck[1]
			userid, _ := strconv.ParseInt(idstr, 10, 0)
			if userid > 0 {
				var user modle.User
				user.Id = userid
				if user.Read() == nil && password == modle.Md5("whomm|"+user.Password) {
					this.userid = user.Id
					this.username = user.Name
					this.nickname = user.Name

					var ug modle.UserGroup
					ug.Id = user.Ugid

					if ug.Read() == nil {
						this.mygroup = ug
						log.Println(ug.Name)
					}

				}
			}
		}
		if this.userid == 0 {
			this.Redirect(this.loginurl, 302)
		}
	}
}

//后台信息提示页
func (this *baseController) showmsg(msg ...string) {
	if len(msg) == 1 {
		msg = append(msg, "javascript:history.back(-1);")
	}
	this.Data["adminid"] = this.userid
	this.Data["adminname"] = this.username
	this.Data["msg"] = msg[0]
	this.Data["redirect"] = msg[1]
	this.TplName = "showmsg.html"
	this.Render()
	this.StopRun()
}

///client ip
func (this *baseController) getClientIp() string {
	return this.Ctx.Input.IP()
}

///获取当前时间
func (this *baseController) getTime() time.Time {
	return time.Now()
}

//是否post提交
func (this *baseController) IsPost() bool {
	return this.Ctx.Request.Method == "POST"
}

//format time
func (this *baseController) FormatTime(t time.Time, format string) string {
	res := strings.Replace(format, "MM", t.Format("01"), -1)
	res = strings.Replace(res, "M", t.Format("1"), -1)
	res = strings.Replace(res, "DD", t.Format("02"), -1)
	res = strings.Replace(res, "D", t.Format("2"), -1)
	res = strings.Replace(res, "YYYY", t.Format("2006"), -1)
	res = strings.Replace(res, "YY", t.Format("06"), -1)
	res = strings.Replace(res, "HH", fmt.Sprintf("%02d", t.Hour()), -1)
	res = strings.Replace(res, "H", fmt.Sprintf("%d", t.Hour()), -1)
	res = strings.Replace(res, "hh", t.Format("03"), -1)
	res = strings.Replace(res, "h", t.Format("3"), -1)
	res = strings.Replace(res, "mm", t.Format("04"), -1)
	res = strings.Replace(res, "m", t.Format("4"), -1)
	res = strings.Replace(res, "ss", t.Format("05"), -1)
	res = strings.Replace(res, "s", t.Format("5"), -1)
	return res
}

func (this *baseController) GetGUID() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(crand.Reader, b); err != nil {
		return ""
	}
	return modle.Md5(base64.URLEncoding.EncodeToString(b))
}

//显示分页链接
func (this *baseController) PageList(pagesize, page, recordcount int64, first bool, path string, param string) (pager string) {
	if recordcount == 0 {
		return ""
	}

	var pagecount int64
	pagecount = 0

	if recordcount%pagesize == 0 {
		pagecount = recordcount / pagesize
	} else {
		pagecount = (recordcount / pagesize) + 1
	}

	pager = "<span>" + strconv.FormatInt(page, 10) + "/" + strconv.FormatInt(pagecount, 10) + "</span>"

	if pagecount < 2 {
		return "<span>共1页</span>"
	}

	pager = pager + "<a href=\"" + path + "/\\?" + param + "\">第一页</a>"

	if page > 1 {
		if page == 2 {
			pager = pager + "<a href=\"" + path + "/\\?" + param + "\">上一页</a>"
		} else {
			pager = pager + "<a href=\"" + path + "/" + strconv.FormatInt(page-1, 10) + "/\\?" + param + "\" >上一页</a>"
		}
	} else {
		pager = pager + "<a href=\"" + path + "/\\?" + param + "\">上一页</a>"
	}

	if page < pagecount {
		pager = pager + "<a href=\"" + path + "/" + strconv.FormatInt(page+1, 10) + "/\\?" + param + "\" class=\"next\">下一页</a>"
	} else {
		pager = pager + "<a href=\"" + path + "/" + strconv.FormatInt(pagecount, 10) + "/\\?" + param + "\"  class=\"next\">下一页</a>"
	}

	pager = pager + "<a href=\"" + path + "/" + strconv.FormatInt(pagecount, 10) + "/\\?" + param + "\"  class=\"next\">最后一页</a>"

	return pager

}

//显示分页链接
func (this *baseController) PageArr(pagesize, page, recordcount int64, first bool, path string, param string) (pager map[string]string) {

	ret := map[string]string{
		"Begin":    "#",
		"Previous": "#",
		"Next":     "#",
		"End":      "#",
	}

	if recordcount == 0 {
		return ret
	}

	var pagecount int64
	pagecount = 0

	if recordcount%pagesize == 0 {
		pagecount = recordcount / pagesize
	} else {
		pagecount = (recordcount / pagesize) + 1
	}

	if pagecount < 2 {

		return ret
	}

	ret["Begin"] = path + "/?" + param

	if page > 1 {
		if page == 2 {
			ret["Previous"] = path + "/?" + param
		} else {
			ret["Previous"] = path + "/" + strconv.FormatInt(page-1, 10) + "/?" + param
		}
	} else {
		ret["Previous"] = path + "/?" + param
	}

	if page < pagecount {

		ret["Next"] = path + "/" + strconv.FormatInt(page+1, 10) + "/?" + param
	} else {

		ret["Next"] = path + "/" + strconv.FormatInt(pagecount, 10) + "/?" + param
	}
	ret["End"] = path + "/" + strconv.FormatInt(pagecount, 10) + "/?" + param

	return ret

}

//显示分页链接
func (this *baseController) Pageinfo(pagesize, page, recordcount int64, first bool, path string) (pager string) {
	if recordcount == 0 {
		return ""
	}
	var pagecount int64
	pagecount = 0
	if recordcount%pagesize == 0 {
		pagecount = recordcount / pagesize
	} else {
		pagecount = (recordcount / pagesize) + 1
	}
	return "Showing " + strconv.FormatInt(pagesize*(page-1)+1, 10) + " to " + strconv.FormatInt(pagesize*page, 10) + " of " + strconv.FormatInt(recordcount, 10) + " entries (" + strconv.FormatInt(page, 10) + "/" + strconv.FormatInt(pagecount, 10) + ")"

}
