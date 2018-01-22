package modle

import (
	"github.com/astaxie/beego/orm"
)

type UserGroup struct {
	Id    int64
	Name  string `orm:"size(50)"`
	Level string
}

func init() {
	// register model
	orm.RegisterModelWithPrefix("t_", new(UserGroup))
}

func (m *UserGroup) TableName() string {
	return "user_group"
}

func (m *UserGroup) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}

func (m *UserGroup) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *UserGroup) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *UserGroup) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *UserGroup) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *UserGroup) GetList(page int64, pagesize int64) []*UserGroup {
	var info UserGroup
	list := make([]*UserGroup, 0)
	info.Query().OrderBy("id").Limit(pagesize, page*pagesize).All(&list)
	return list
}
