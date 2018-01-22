package modle

import (
	"github.com/astaxie/beego/orm"
)

type User struct {
	Id       int64
	Name     string `orm:"size(30)"`
	Password string `orm:"size(32)"`
	Ugid     int64
}

func init() {
	// register model
	orm.RegisterModelWithPrefix("t_", new(User))
}

func (m *User) TableName() string {
	return "user"
}

func (m *User) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *User) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *User) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
