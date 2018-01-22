package modle

import (
	"encoding/json"
	"errors"

	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	hprose "github.com/hprose/hprose-go"
)

//任务运行状态
const (
	TS_PENGDING = iota //阻塞
	TS_READY           //就绪
	TS_RUN             //运行
	TS_SUCCESS         //成功
	TS_FAIL            //失败
	TS_KILLING         //杀死中
	TS_KILLED          //已杀死
	TS_CALLING         //通知中（通知运行）
	TS_CALLFAIL        //通知失败
)

type Instance struct {
	Id            int64 `orm:auto`
	Tid           int64
	Runtime       int64
	Tasktime      int64
	Status        int
	Createby      int64
	Time_create   int64
	Time_ready    int64
	Time_run      int64
	Time_success  int64
	Time_fail     int64
	Time_callfail int64
	Time_killing  int64
	Time_killed   int64
	Relaystatus   string
	Calltime      int64
}

const (
	INSTANCETABE = "instance"
)

var (
	ErrRecall = errors.New("recall")
)

func init() {
	// register model
	orm.RegisterModelWithPrefix(TABLEPRE, new(Instance))
}

//orm的获取表名称
func (t *Instance) TableName() string {
	return INSTANCETABE
}

//创建
func (Instance) Create(runtime int64, tasktime int64, tid int64, createby int64) (*Instance, error) {
	its := Instance{Tid: tid, Runtime: runtime, Tasktime: tasktime, Createby: createby, Status: 0, Time_create: time.Now().Unix()}
	o := orm.NewOrm()
	id, err := o.Insert(&its)
	if err == nil {
		its.Id = id
	} else {
		return nil, err
	}
	return &its, nil

}

//更新为成功
func (its *Instance) UpdateRuning() error {

	o := orm.NewOrm()
	err := o.Begin()

	res, err := o.Raw("update t_instance set status=? where id = ? and  status in (?)", TS_RUN, its.Id, TS_CALLFAIL).Exec()

	if err != nil {
		err = o.Rollback()
		return err
	} else {
		num, _ := res.RowsAffected()
		err = o.Commit()
		if err != nil {
			return errors.New("db transaction error")
		}
		if num == 1 {
			return nil
		} else {
			return errors.New("change instance status error")
		}

	}

}

//更新为ready
func (its *Instance) UpdateReady() error {

	o := orm.NewOrm()
	err := o.Begin()

	res, err := o.Raw("update t_instance set status=? where id = ? and  status in (?)", TS_READY, its.Id, TS_PENGDING).Exec()

	if err != nil {
		err = o.Rollback()
		return err
	} else {
		num, _ := res.RowsAffected()
		err = o.Commit()
		if err != nil {
			return errors.New("db transaction error")
		}
		if num == 1 {
			return nil
		} else {
			//有的任务已经更新为ready,影响行数可能为0不用返回error
			return nil
		}

	}

}

//更新为通知失败
func (its *Instance) UpdateCallfail() error {

	o := orm.NewOrm()
	err := o.Begin()
	// 事务处理过程
	res, err := o.Raw("update t_instance set status=? where id = ? and  status in (?)", TS_CALLFAIL, its.Tid, TS_CALLING, TS_READY).Exec()
	// 此过程中的所有使用 o Ormer 对象的查询都在事务处理范围内
	if err == nil {
		num, _ := res.RowsAffected()
		if num == 1 {
			_, err = o.Raw("delete from t_runing_tag where tid=? and tasktime=?", its.Tid, its.Tasktime).Exec()
		} else {
			err = o.Rollback()
			if err == nil {
				return errors.New("update instance status error")
			}
			return err
		}

	}
	if err != nil {
		err = o.Rollback()

	} else {
		err = o.Commit()
	}
	return err

}

//更新为kiled状态
func (its *Instance) UpdateKilled() (int64, error) {

	o := orm.NewOrm()
	err := o.Begin()
	res, err := o.Raw("update t_instance set status=? where id = ? and  status in (?,?)", TS_KILLED, its.Id, TS_PENGDING, TS_READY).Exec()
	var num int64
	if err == nil {
		num, _ := res.RowsAffected()
		if num == 1 {
			_, err = o.Raw("delete from t_runing_tag where tid=? and tasktime=?", its.Tid, its.Tasktime).Exec()
		} else {
			err = o.Rollback()
			if err == nil {
				return 0, errors.New("update instance status error")
			}
			return 0, err
		}

	}
	if err != nil {
		err = o.Rollback()

	} else {
		err = o.Commit()
	}
	return num, err

}

//更新为killing状态
func (its *Instance) UpdateKilling() (int64, error) {

	o := orm.NewOrm()
	err := o.Begin()

	res, err := o.Raw("update t_instance set status=? where id = ? and  status in (?)", TS_KILLING, its.Id, TS_RUN).Exec()
	var num int64
	if err == nil {
		num, _ = res.RowsAffected()
	}
	if err != nil {
		err = o.Rollback()

	} else {
		err = o.Commit()
	}
	return num, err

}

//获取instance
func (its *Instance) IsKilled() (bool, error) {

	o := orm.NewOrm()
	t := Instance{Id: its.Id}

	err := o.Read(&t)
	if err == nil {
		return t.Status == TS_KILLED, err
	} else if err == orm.ErrNoRows {
	} else if err == orm.ErrMissPK {

	} else {

	}
	return false, err

}

//获取instance
func (its *Instance) GetId(id int64) *Instance {

	o := orm.NewOrm()
	t := Instance{Id: id}

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

//通过id获取taskinfo
func (its *Instance) GetTaskinfobyId(id int64) *Task {

	o := orm.NewOrm()
	t := Instance{Id: id}

	err := o.Read(&t)
	if err == nil {
		t := (Task{}).GetTaskbyid(t.Tid)
		return t
	} else if err == orm.ErrNoRows {
		return nil
	} else if err == orm.ErrMissPK {
		return nil
	} else {
		return nil
	}

}

//增加通知次数
func (its *Instance) AddCallTime() (int64, error) {

	o := orm.NewOrm()
	err := o.Begin()
	var calltime int64 = 0
	// 事务处理过程
	res, err := o.Raw("update t_instance set calltime=calltime+1 where id = ? and  status in (?)", its.Id, TS_READY).Exec()
	// 此过程中的所有使用 o Ormer 对象的查询都在事务处理范围内
	if err == nil {
		num, _ := res.RowsAffected()
		if num == 1 {

			t := Instance{Id: its.Id}
			err := o.Read(&t)
			if err != nil {
				//todo:判断一下err
			}
			calltime = t.Calltime

		} else {
			err = o.Rollback()
			if err == nil {
				return 0, errors.New("update instance calltime error")
			}
			return 0, err
		}

	}
	if err != nil {
		err = o.Rollback()

	} else {
		err = o.Commit()
	}
	return calltime, err

}

//再次更新为pendding状态
func (its *Instance) ReUpdatePennding(relaystr string) (int64, error) {

	o := orm.NewOrm()
	err := o.Begin()
	res, err := o.Raw("update t_instance set status=?, relaystatus=?, calltime=0 where id = ? and  status in (?,?)", TS_PENGDING, relaystr, its.Id, TS_PENGDING, TS_READY).Exec()
	var num int64
	if err == nil {
		num, _ := res.RowsAffected()
		return num, err
	}
	if err != nil {
		err = o.Rollback()

	} else {
		err = o.Commit()
	}
	return num, err

}

//通知客户端运行
func (its *Instance) Calltorun(taskinfo Task) error {

	o := orm.NewOrm()
	err := o.Begin()

	//定义一个运行异常处理函数，用于通知上层重跑的
	defer func() {
		if r := recover(); r != nil {
			logs.Error(r)
			err = o.Rollback()
			if err != nil {
				//panic(err)
				//这个时候不知道怎么处理了
			}
			//传递到调用函数方
			panic(errors.New("recall"))

		}
	}()

	//更新实例CALLING状态
	res, err := o.Raw("update t_instance set status=? where id = ? and  status in (?)", TS_CALLING, its.Id, TS_READY).Exec()

	if err != nil {
		//panic(errors.New(" ready  => calling sql error "))
		goto CALLFAIL
	} else {
		num, _ := res.RowsAffected()
		if num != 1 {
			//panic(errors.New("ready to calling affect row not one"))
			goto CALLFAIL
		} else {
			//CALLING DONE
			var f interface{}
			if err := json.Unmarshal([]byte(taskinfo.Extra), &f); err != nil {
				//没有通知通配置
				//panic(errors.New("no call worker configuie error"))
				goto CALLFAIL
			} else {
				extra := f.(map[string]interface{})
				if extra["workrpc"] == nil {
					//work的配置有问题
					//panic(errors.New("worker configuie error"))
					goto CALLFAIL
				} else {
					//hpros maybe throw panic
					client := hprose.NewClient(extra["workrpc"].(string))
					var ro *clientStub
					client.UseService(&ro)

					extra["taskname"] = taskinfo.Name
					jtinfo, err := json.Marshal(extra)
					if err != nil {
						//panic(errors.New("before call work , format info error callerror"))
						goto CALLFAIL
					}

					if ro.Run(its.Tasktime, its.Tid, its.Id, string(jtinfo)) {
						logs.Debug("call work success set job success")
						//通知成功了
						res, err := o.Raw("update t_instance set status=? where id = ? and  status in (?)", TS_RUN, its.Id, TS_CALLING).Exec()
						if err != nil {
							//panic(errors.New("calling to run success but update status error"))
							//这个时候不知道怎么处理了
							//todo
						}
						r, _ := res.RowsAffected()
						logs.Debug("RowsAffected:%d", r)

					} else {
						//panic(errors.New("call worker error"))
						goto RECALL
					}

					err = o.Commit()
					if err != nil {
						//errors.New("db transaction error")
						//这个时候不知道怎么处理了
						//todo
					}
					return nil

				}

			}

		}

	}

	//这个任务还可以重试一下
RECALL:
	panic(ErrRecall)

	//标记任务失败
CALLFAIL:

	res, err = o.Raw("update t_instance set status=? where id = ? and  status in (?,?)", TS_CALLFAIL, its.Id, TS_CALLING, TS_READY).Exec()
	if err != nil {
		//panic(err)
		//这个时候不知道怎么处理了
		//todo
	}
	_, err = o.Raw("delete from t_runing_tag where tid=? and tasktime=?", its.Tid, its.Tasktime).Exec()
	if err != nil {
		//panic(err)
		//这个时候不知道怎么处理了
		//todo
	}

	err = o.Commit()
	if err != nil {
		//panic(err)
		//这个时候不知道怎么处理了
		//todo
	}
	return nil

}

func (its *Instance) Scuccess() (int64, error) {

	o := orm.NewOrm()
	err := o.Begin()

	res, err := o.Raw("update t_instance set status=? where id = ? and  status in (?,?)", TS_SUCCESS, its.Id, TS_RUN, TS_KILLING).Exec()
	var num int64
	if err == nil {
		num, err = res.RowsAffected()
		if err == nil && num == 1 {
			res, err = o.Raw("update t_runing_tag set status=1 where tid=? and tasktime=?", its.Tid, its.Tasktime).Exec()
			num, err = res.RowsAffected()
		}

	}

	if err != nil {
		err = o.Rollback()

	} else {
		err = o.Commit()
	}
	return num, err

}

func (its *Instance) Fail() (int64, error) {

	o := orm.NewOrm()
	err := o.Begin()

	res, err := o.Raw("update t_instance set status=? where id = ? and  status in (?,?)", TS_FAIL, its.Id, TS_RUN, TS_KILLING).Exec()
	var num int64
	if err == nil {
		num, err = res.RowsAffected()
		if err == nil && num == 1 {
			res, err = o.Raw("delete from t_runing_tag  where tid=? and tasktime=?", its.Tid, its.Tasktime).Exec()
			num, err = res.RowsAffected()
		}

	}

	if err != nil {
		err = o.Rollback()

	} else {
		err = o.Commit()
	}
	return num, err

}

func (its *Instance) Killed() (int64, error) {

	o := orm.NewOrm()
	err := o.Begin()

	res, err := o.Raw("update t_instance set status=? where id = ? and  status in (?)", TS_KILLED, its.Id, TS_KILLING).Exec()
	var num int64
	if err == nil {
		num, err = res.RowsAffected()
		if err == nil && num == 1 {
			res, err = o.Raw("delete from t_runing_tag  where tid=? and tasktime=?", its.Tid, its.Tasktime).Exec()
			num, err = res.RowsAffected()
		}
	}

	if err != nil {
		err = o.Rollback()

	} else {
		err = o.Commit()
	}
	return num, err

}

func (m *Instance) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
