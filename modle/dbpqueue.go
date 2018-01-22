//数据库优先级队列
package modle

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"errors"

	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

const (
	GENERATETABE = "pq_generate"
	RUNTETABE    = "pq_run"
)

var (
	ErrPush = errors.New("pqueue create one error")
	ErrPop  = errors.New("pqueue pop one error")
)

func init() {
	gob.Register(DbstoreItem{})
	gob.Register(Task{})
	gob.Register(Instance{})
}

//数据库字段存储数据结构
type DbstoreItem struct {
	Data interface{}
}

//对外接口
type IPQueue interface {
	Push(task interface{}, priority int64) (int64, error)
	Pop() (interface{}, error)
}

//获取表格接口
type Tableinfo interface {
	GetTableName() string
}

//队列item
type DbPQueue struct {
	Id       int64 `orm:auto`
	Taskins  string
	Priority int64
	Done     int
	tinfo    Tableinfo
}

//取表名称
func (t *DbPQueue) TableName() string {
	return t.tinfo.GetTableName()
}

//push
func (t *DbPQueue) Push(task interface{}, priority int64) (int64, error) {
	gob.Register(task)

	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(DbstoreItem{Data: task})

	if err != nil {
		return 0, ErrPush
	}

	o := orm.NewOrm()
	res, err := o.Raw("insert into   "+TABLEPRE+t.TableName()+" set done = 0, taskins=?, priority=? ", hex.EncodeToString(buf.Bytes()), priority).Exec()
	if err != nil {
		logs.Error("db queue insert error %s", err)
		return 0, ErrPush
	} else {
		id, err := res.LastInsertId()
		if err == nil {
			return id, nil
		} else {
			logs.Error("db queue insert error  error %s", err)
			return 0, ErrPush
		}
	}
}

//pop
func (t *DbPQueue) Pop() (interface{}, error) {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return nil, ErrPop
	}
	var itm DbPQueue

	err = o.Raw("select * from  "+TABLEPRE+t.TableName()+" where priority<=? and done=0 order by priority desc limit 1", time.Now().Unix()).QueryRow(&itm)
	if err != nil {
		return nil, ErrPop
	}

	res, err := o.Raw("update  "+TABLEPRE+t.TableName()+" set done = 1  where id=?", itm.Id).Exec()
	if err != nil {
		err = o.Rollback()
		logs.Error("db queue get update table error %s", err)
		return nil, ErrPop
	} else {
		num, _ := res.RowsAffected()
		err = o.Commit()
		if err != nil {
			logs.Error("db queue get transaction error %s", err)
			return nil, ErrPop
		}
		if num == 1 {

			logs.Info("pop from table %d", itm.Id)
			bs, _ := hex.DecodeString(itm.Taskins)
			buf := bytes.NewBuffer(bs)
			var item DbstoreItem
			dec := gob.NewDecoder(buf)
			err := dec.Decode(&item)
			if err == nil {
				return item.Data, nil
			}
			logs.Info("pop decode error %s", err)
			return nil, ErrPop
		} else {
			//有的任务已经更新为ready,影响行数可能为0不用返回error
			return nil, ErrPop
		}

	}
}

//任务生成队列
type GenTableinfo struct {
}

func (p *GenTableinfo) GetTableName() string {
	return GENERATETABE
}

//任务运行队列
type RunTableinfo struct {
}

func (p *RunTableinfo) GetTableName() string {
	return RUNTETABE
}

//创建一个数据库队列
func NewDBPQueue(t string) IPQueue {
	if t == "gen" {
		return &DbPQueue{tinfo: &GenTableinfo{}}
	} else {
		return &DbPQueue{tinfo: &RunTableinfo{}}

	}
}
