package modle

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"time"

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
	res, err := o.Raw("insert into   "+TABLEPRE+t.TableName()+" set done = 0, taskins=?, priority=? ", buf.String(), priority).Exec()
	if err != nil {
		log.Println("db queue insert error ", err)
		return 0, ErrPush
	} else {
		id, err := res.LastInsertId()
		if err == nil {
			return id, nil
		} else {
			log.Println("db queue insert error  error", err)
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
		log.Println("db queue get update table error", err)
		return nil, ErrPop
	} else {
		num, _ := res.RowsAffected()
		err = o.Commit()
		if err != nil {
			log.Println("db queue get transaction error", err)
			return nil, ErrPop
		}
		if num == 1 {

			buf := bytes.NewBufferString(itm.Taskins)
			var item DbstoreItem
			dec := gob.NewDecoder(buf)
			err := dec.Decode(&item)
			if err != nil {
				return item.Data, ErrPush
			}
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
