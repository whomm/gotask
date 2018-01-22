package modle

import (
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type Runingtag struct {
	Id       int64 `orm:auto`
	Tid      int64
	Tasktime int64
	Status   int
}

const (
	RUNGINGTABE = "runing_tag"
)

func init() {
	// register model
	orm.RegisterModelWithPrefix(TABLEPRE, new(Runingtag))
}

//orm的获取表名称
func (t *Runingtag) TableName() string {
	return RUNGINGTABE
}

//创建任务运行标记,已确保只有一个处于运行状态
func (Runingtag) Create(runtime int64, tasktime int64, tid int64, createby int64) (*Instance, error) {
	tg := Runingtag{Tid: tid, Tasktime: tasktime, Status: 0}
	its := Instance{Tid: tid, Runtime: runtime, Tasktime: tasktime, Createby: createby, Status: 0, Time_create: time.Now().Unix()}

	o := orm.NewOrm()
	err := o.Begin()
	// 事务处理过程
	//删除处于成功状态的任务
	_, err = o.Raw("delete from t_runing_tag where tasktime = ? and tid=? and status>0", tasktime, tid).Exec()
	if err == nil {

		//添加处于活动状态的任务
		id, err := o.Insert(&tg)
		if err == nil {
			tg.Id = id

			//开始创建任务实例

			id, err := o.Insert(&its)
			if err == nil {
				its.Id = id
			}
		}
	}

	// 此过程中的所有使用 o Ormer 对象的查询都在事务处理范围内
	if err != nil {
		err = o.Rollback()
		return nil, err
	} else {
		err = o.Commit()
		return &its, err
	}

}

//检查任务依赖
func (Runingtag) Checkrely(runtime int64, relay string) (bool, string) {

	//todo
	return true, ""

}
