package modle

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron"
)

type relayitem struct {
	start  string
	length string
	id     int64
}

type relays struct {
	rl []relayitem
}

type Task struct {
	Id          int64  `orm:auto`
	Uid         int64  `valid:"Required"`
	Ugid        int64  `valid:"Required"`
	Tgid        int64  `valid:"Required"`
	Name        string `valid:"Required;MinSize(1)"`
	Crontab     string
	Nextrun     uint
	Pendingtime uint
	Starttime   uint
	Endtime     uint
	Extra       string
	Invalid     int
	Relay       string `valid:"Required"`
}

const (
	TABLEPRE = "t_"
	TASKTABE = "task"
)

func init() {
	// register model
	orm.RegisterModelWithPrefix(TABLEPRE, new(Task))
}

//orm的获取表名称
func (t *Task) TableName() string {
	return TASKTABE
}

//更新下次运行时间
func (t *Task) UpdateNextRuntime(timenow int64) (bool, error) {

	sc, err := cron.Parse(t.Crontab)
	if err == nil {
		if timenow > 0 {
			t.Nextrun = uint(sc.Next(time.Unix(timenow, 0)).Unix()) //+ t.Pendingtime
		} else {
			t.Nextrun = uint(sc.Next(time.Now()).Unix()) //+ t.Pendingtime
		}

		log.Println(t.Crontab)
		log.Println(sc.Next(time.Now()))

		return true, nil
	} else {
		return false, errors.New("crontab parse error")
	}

}

func taskidvalide(ids []int64) bool {

	o := orm.NewOrm()
	cnt, err := o.QueryTable("t_task").Filter("id__in", ids).Filter("invalid__gt", 0).Count()
	if err != nil {
		return false
	}
	if cnt != int64(len(ids)) {
		return false
	}
	return true
}

// 如果你的 struct 实现了接口 validation.ValidFormer
// 当 StructTag 中的测试都成功时，将会执行 Valid 函数进行自定义验证
func (t *Task) Valid(v *validation.Validation) {
	if strings.Index(t.Name, "admin") != -1 {
		// 通过 SetError 设置 Name 的错误信息，HasErrors 将会返回 true
		v.SetError("Name", "名称里不能含有 admin")
	}

	if t.Id < 0 {
		v.SetError("Id", "必须大于等于0")
	}

	if t.Crontab != "" {
		_, err := cron.Parse(t.Crontab)
		if err != nil {
			v.SetError("Corntab", "格式不正确")
		}
	}
	//校验依赖格式是否正确 {id:120,start:-1day,length:1day}
	//id所属的用户组
	var rs relays
	err := json.Unmarshal([]byte(t.Relay), &rs)
	//log.Println(err)
	if err != nil {
		v.SetError("Relay", "格式不正确")
	}
	if len(rs.rl) > 0 {
		ids := []int64{}
		for _, i := range rs.rl {
			ids = append(ids, i.id)
		}
		if !taskidvalide(ids) {
			v.SetError("Relay", "id不合法")
		}
	}
}

//获取任务
func (Task) GetTaskbyid(id int64) *Task {
	o := orm.NewOrm()
	t := Task{Id: id}

	err := o.Read(&t)

	if err == nil {
		return &t
	} else if err == orm.ErrNoRows {
		return nil
	} else if err == orm.ErrMissPK {
		return nil
	} else {
		return nil
	}
}

func (m *Task) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
