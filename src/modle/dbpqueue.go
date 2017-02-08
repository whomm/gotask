package modle

import (
	"encoding/json"
	"log"
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	GENERATETABE = "pq_generate"
	RUNTETABE    = "pq_run"
)

type IPQtable interface {
	Create(taskins string, priority int64, done int) (int64, error)
	GetOne() string
}

type PQGenerate struct {
	Id       int64 `orm:auto`
	Taskins  string
	Priority int64
	Done     int
}

//orm的获取表名称
func (t *PQGenerate) TableName() string {
	return GENERATETABE
}

//创建
func (t *PQGenerate) Create(taskins string, priority int64, done int) (int64, error) {
	its := PQGenerate{Taskins: taskins, Priority: priority, Done: done}
	o := orm.NewOrm()
	id, err := o.Insert(&its)
	if err == nil {
		its.Id = id
	} else {
		return 0, err
	}
	return its.Id, nil
}

func (pq *PQGenerate) GetOne() string {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return ""
	}
	var itm PQGenerate

	err = o.Raw("select * from  "+TABLEPRE+pq.TableName()+" where priority<=? and done=0 order by priority desc limit 1", time.Now().Unix()).QueryRow(&itm)
	if err != nil {
		return ""
	}

	res, err := o.Raw("update  "+TABLEPRE+pq.TableName()+" set done = 1  where id=?", itm.Id).Exec()
	if err != nil {
		err = o.Rollback()
		log.Println("db queue get update table error", err)
		return ""
	} else {
		num, _ := res.RowsAffected()
		err = o.Commit()
		if err != nil {
			log.Println("db queue get transaction error", err)
			return ""
		}
		if num == 1 {
			return itm.Taskins
		} else {
			//有的任务已经更新为ready,影响行数可能为0不用返回error
			return ""
		}

	}
}

type PQRun struct {
	Id       int64 `orm:auto`
	Taskins  string
	Priority int64
	Done     int
}

//orm的获取表名称
func (t *PQRun) TableName() string {
	return RUNTETABE
}

//创建
func (t *PQRun) Create(taskins string, priority int64, done int) (int64, error) {
	its := PQRun{Taskins: taskins, Priority: priority, Done: done}
	o := orm.NewOrm()
	id, err := o.Insert(&its)
	if err == nil {
		its.Id = id
	} else {
		return 0, err
	}
	return its.Id, nil
}

func (pq *PQRun) GetOne() string {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return ""
	}
	var itm PQRun

	err = o.Raw("select * from  "+TABLEPRE+pq.TableName()+" where priority<=? and done=0 order by priority desc limit 1", time.Now().Unix()).QueryRow(&itm)
	if err != nil {
		return ""
	}

	res, err := o.Raw("update  "+TABLEPRE+pq.TableName()+" set done = 1  where id=?", itm.Id).Exec()
	if err != nil {
		err = o.Rollback()
		log.Println("db queue get update table error", err)
		return ""
	} else {
		num, _ := res.RowsAffected()
		err = o.Commit()
		if err != nil {
			log.Println("db queue get transaction error", err)
			return ""
		}
		if num == 1 {
			return itm.Taskins
		} else {
			//有的任务已经更新为ready,影响行数可能为0不用返回error
			return ""
		}

	}
}

func init() {
	// register model
	orm.RegisterModelWithPrefix(TABLEPRE, new(PQGenerate))
	orm.RegisterModelWithPrefix(TABLEPRE, new(PQRun))
}

type DBPQueue struct {
	//sync.RWMutex
	ptable IPQtable
}

//创建一个数据库队列
func NewDBPQueue(t string) *DBPQueue {
	if t == "gen" {
		x := DBPQueue{}
		x.ptable = &PQGenerate{}
		return &x
	} else {
		x := DBPQueue{}
		x.ptable = &PQRun{}
		return &x

	}
}

func (pq *DBPQueue) Pop() (interface{}, int64) {

	str := (pq.ptable).(IPQtable).GetOne()
	if str != "" {
		var its map[string]interface{}
		err := json.Unmarshal([]byte(str), &its)
		if err != nil {
			log.Println("error:", err)
		}
		return &its, 0
	}
	return nil, 0
}

// Push the value item into the priority queue with provided priority.
func (pq *DBPQueue) Push(value interface{}, priority int64) error {

	b, err := json.Marshal(value)
	if err != nil {
		log.Println("pq queue error:", err)
		return err
	}

	_, err = (pq.ptable).(IPQtable).Create(string(b), priority, 0)
	if err != nil {
		log.Println("pq queue insert data error", err)
		return err
	}

	return nil

}
