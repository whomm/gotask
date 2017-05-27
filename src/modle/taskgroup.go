/*
   任务组:用来管理任务的分类,属于某个用户组，用户组可以创建多个任务组方便任务的管理
*/
package modle

import (
	"github.com/astaxie/beego/orm"
)

type TaskGroup struct {
	Id    int64
	Name  string `orm:"size(50)"`
	Extra string
	Uid   int64
	Ugid  int64
}

func init() {
	// register model
	orm.RegisterModelWithPrefix("t_", new(TaskGroup))
}

func (m *TaskGroup) TableName() string {
	return "task_group"
}

func (m *TaskGroup) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}

func (m *TaskGroup) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *TaskGroup) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *TaskGroup) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *TaskGroup) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *TaskGroup) GetList(page int64, pagesize int64) []*TaskGroup {
	var info TaskGroup
	list := make([]*TaskGroup, 0)
	info.Query().OrderBy("id").Limit(pagesize, page*pagesize).All(&list)
	return list
}

func (m *TaskGroup) GetAllByUgid(ugid int64) []*TaskGroup {

	var list []*TaskGroup
	query := m.Query()
	query = query.Filter("ugid__eq", ugid)
	count, _ := query.Count()
	if count > 0 {
		query.OrderBy("-Id").All(&list)
	}
	return list
}
