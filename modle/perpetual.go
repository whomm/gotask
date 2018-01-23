package modle

import (
	"encoding/json"
	"errors"

	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	hprose "github.com/hprose/hprose-go"
	"sync/atomic"
)


func init() {


	//设置最大空闲连接
	//设置最大数据库连接 (go >= 1.2)
	maxIdle := 30
	maxConn := 30
	// set default database
	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("mysqlurls"), maxIdle, maxConn)
	orm.Debug = false
}

type PManage struct {
	NowStop int32
}

//获取任务
func (p *PManage) GetTask(id int64) (Task, error) {
	o := orm.NewOrm()
	t := Task{Id: id}

	err := o.Read(&t)

	if err == orm.ErrNoRows {
		return t, errors.New("can not find")
	} else if err == orm.ErrMissPK {
		return t, errors.New("can not find PK")
	} else {
		return t, nil
	}
}

//创建任务
func (p *PManage) CreateTask(t Task) (int64, error) {

	//校验任务是否合法
	valid := validation.Validation{}
	b, err := valid.Valid(&t)
	if err != nil {
		return 0, err
	}

	if !b {
		einfo := ""
		for _, err := range valid.Errors {
			logs.Info("%s, %s", err.Key, err.Message)
			einfo = einfo + err.Key + err.Message + "; "
		}
		return 0, errors.New(einfo)
	}
	t.Nextrun = 0
	t.Id = 0

	//直接数据库创建任务返回（任务无效，开始时间大约结束时间，结束时间大于当前时间）
	if t.Invalid == 1 || t.Starttime >= t.Endtime || t.Endtime < uint(time.Now().Unix()) {
		o := orm.NewOrm()
		return o.Insert(&t)
	}

	//更新下一次运行时间
	_, err = t.UpdateNextRuntime(0)
	if err != nil {
		return 0, errors.New("compute next run time error before create task ")
	} else {

		//数据库创建并添加到实例生成队列
		o := orm.NewOrm()
		id, err := o.Insert(&t)
		if err != nil {
			return id, err
		}

		if t.Nextrun < t.Endtime {
			p.PushGenerater(t)
			logs.Debug("add to queue[create]")
			return id, nil
		} else {
			return id, nil
		}
	}

}

//更新任务
func (p *PManage) UpdateTask(t Task, field []string) (bool, error) {
	o := orm.NewOrm()
	if t.Id == 0 {
		return false, errors.New("id is 0")
	} else {

		t.Nextrun = 0
		//直接数据库更新任务返回（任务无效，开始时间大约结束时间，结束时间大于当前时间），任务生成器会在后续校验的时候自动去掉该任务
		if t.Invalid == 1 || t.Starttime >= t.Endtime || t.Endtime < uint(time.Now().Unix()) {
			_, err := o.Update(&t, field...)
			if err != nil {
				return false, err
			}
			return true, nil
		}

		_, err := t.UpdateNextRuntime(0)

		if err != nil {
			return false, err
		}

		//更新数据库并添加到实例生成队列
		_, err = o.Update(&t, field...)
		if err != nil {
			return false, err
		}

		if t.Nextrun < t.Endtime {
			p.PushGenerater(t)
			logs.Debug("add to queue[update]")
			return true, nil
		}
		return true, nil
	}
	return false, errors.New("unkonwn error")
}

//删除任务  数据库标记任务失效，并非真正删除（正在运行中的任务没有处理）
func (p *PManage) RemoveTask(id int64) (int64, error) {
	o := orm.NewOrm()
	res, err := o.Raw("UPDATE t_task SET Invalid = 1 where id=?", id).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		return num, err
	}
	return 0, err

}

//任务实例生成队列
var TaskGenerater IPQueue = NewDBPQueue("gen")

//任务实例运行调度队列
var TaskRuner IPQueue = NewDBPQueue("run")

//添加任务到生成器队列
func (p *PManage) PushGenerater(t Task) {
	TaskGenerater.Push(t, int64(t.Nextrun))
}

//任务实例生成
func (p *PManage) GeneraterConsume() {
	
	for atomic.LoadInt32(&p.NowStop) < 0 {
		headValue, _ := TaskGenerater.Pop()
		if headValue == nil {
			logs.Debug("got nil one")
		} else {
			logs.Debug("got one")
			//创建任务实例添加到运行队列
			task, err := p.GetTask(headValue.(Task).Id)
			if err != nil {
				//获取不到任务
				logs.Error("get task err %s", err)

			} else {
				if task.Invalid == 1 ||
					task.Starttime >= task.Endtime ||
					task.Endtime <= headValue.(Task).Nextrun ||
					task.Crontab == "" ||
					task.Crontab != headValue.(Task).Crontab {
					//任务被更新过可以丢弃了
					logs.Debug(headValue.(Task).Name + " [time]" + time.Unix(int64(headValue.(Task).Nextrun), 0).Format("2006-01-02 15:04:05") + " [invalide]")
				} else {
					go p.PushRuner(headValue.(Task))
					logs.Debug(headValue.(Task).Name + " [time]" + time.Unix(int64(headValue.(Task).Nextrun), 0).Format("2006-01-02 15:04:05") + " [push to runing queue]")

					//用上一个任务的运行时间生成下一个运行时间，防止漏掉某个时间点
					task.UpdateNextRuntime(int64(headValue.(Task).Nextrun))
					p.PushGenerater(task)

				}

			}
			//拿到一个赶快看看下一个
			time.Sleep(10000 * time.Millisecond)
			continue
		}
		//拿不到了多休息一会
		time.Sleep(10000 * time.Millisecond)

	}
}

//添加运行队列
func (p *PManage) PushRuner(t Task) {

	//创建任务实例 Runingtag确保活动状态任务唯一
	r, err := (Runingtag{}).Create(int64(t.Nextrun+t.Pendingtime), int64(t.Nextrun), t.Id, 0)
	if err != nil {
		//记录一下错误的类型
		logs.Error("create run tag && instance error %s ", err)

	} else {
		TaskRuner.Push(*r, int64(t.Nextrun+t.Pendingtime))
	}

}

//运行队列的消费者
func (p *PManage) RunerConsume() {
	for atomic.LoadInt32(&p.NowStop) < 0 {
		headValue, _ := TaskRuner.Pop()
		if headValue == nil {
			//continue
		} else {
			logs.Debug("[run]got one to run")
			//获取到需要运行的任务
			task, err := p.GetTask(headValue.(Instance).Tid)
			if err == nil {

				logs.Debug("[run] " + task.Name + " [time+pendding]" + time.Unix(int64(headValue.(Instance).Runtime), 0).Format("2006-01-02 15:04:05") + " [time]" + time.Unix(int64(headValue.(Instance).Tasktime), 0).Format("2006-01-02 15:04:05"))

				go p.checkrelayandrun(headValue.(Instance), task)

			}
			//拿到一个赶快看看下一个
			time.Sleep(10000 * time.Millisecond)
			continue
		}

		time.Sleep(10000 * time.Millisecond)

	}
}

//用户端的通知接口
type clientStub struct {
	Run  func(tasktime int64, taskid int64, taskinstanceid int64, taskinfo string) bool
	Kill func(taskinstanceid int64) bool
}

//通知客户端kill任务
func (p *PManage) notifyclientkill(taskinstanceid int64, taskid int64) error {
	taskinfo := Task{}.GetTaskbyid(taskid)
	if taskinfo != nil {
		var f interface{}
		err := json.Unmarshal([]byte(taskinfo.Extra), &f)
		if err != nil {
			return errors.New("extra jsion error")
		}
		extra := f.(map[string]interface{})
		if extra["workrpc"] == nil {
			return errors.New("not rpc  cannot kill call ")
		}

		defer func() {
			if r := recover(); r != nil {
				logs.Error("call rpc  error %s", r)
				//通知失败了
				//to do retry
			}
		}()
		return func() error {
			client := hprose.NewClient(extra["workrpc"].(string))
			var ro *clientStub
			client.UseService(&ro)

			if ro.Kill(taskinstanceid) {
				logs.Debug("call work to kill job success")

				return nil
			} else {
				logs.Error("call worker to kill job  error")

				return errors.New("call worker to kill job  error")

			}
		}()

		return nil

	}
	return errors.New("get task info error")
}

func (p *PManage) killinstance(id int64) error {
	its := Instance{}
	thisit := its.GetId(id)
	if thisit == nil {
		return errors.New("instance not exsit")
	}
	switch thisit.Status {
	case TS_PENGDING, TS_READY:
		//pending和ready的任务可以直接标记为killed
		rs, err := thisit.UpdateKilled()
		if err == nil && rs == 1 {
			return nil
		}
		//如果未能标记成killed，看看是不是running 需要通知客户端杀死
		if err == nil && rs == 0 {
			rs, err = its.UpdateKilling()
			if err == nil && rs == 1 {
				return p.notifyclientkill(id, thisit.Tid)
			}
		}
		return nil
	case TS_RUN:
		//运行中的状态running 通知客户端杀死
		rs, err := thisit.UpdateKilling()
		if err == nil && rs == 1 {
			return p.notifyclientkill(id, thisit.Tid)
		}
		return nil
	case TS_KILLING:
		//killing 再次通知客户端通知

		return p.notifyclientkill(id, thisit.Tid)

	case TS_KILLED:
		//已经是killed状态就不再处理了直接返回
		return nil
	default:
		return errors.New("instance status error")

	}
	return nil
}

//检查依赖调用执行
func (p *PManage) checkrelayandrun(its Instance, taskinfo Task) {

	//看看任务是不是被杀死,杀死的任务就不用跑了
	killed, err := its.IsKilled()
	if err != nil {
		//todo:判断一下
	}
	if killed {

		logs.Debug("[killed not to recall  client] " + taskinfo.Name + " [time+pendding]" + time.Unix(int64(its.Runtime), 0).Format("2006-01-02 15:04:05") + " [time]" + time.Unix(int64(its.Tasktime), 0).Format("2006-01-02 15:04:05"))

		return
	}

	//通知客户端网络异常处理
	defer func() {
		if r := recover(); r != nil {
			logs.Error("call to run error: %s", r)
			if r.(error).Error() == "recall" {
				//调用异常可以重试
				//增加一下运行时间 5 秒，增加调用次数，重新去排队 100次重试可以标记为callfail
				calltime, err := its.AddCallTime()
				if err == nil && calltime > 100 {
					//任务可以标记为callfail
					its.UpdateCallfail()
					return
				}
				TaskRuner.Push(its, time.Now().Unix()+5)

			}
		}
	}()

	//检查依赖
	norely, relyinfo := (Runingtag{}).Checkrely(its.Runtime, taskinfo.Relay)
	if norely {

		//更新ready
		err := its.UpdateReady()
		if err != nil {
			logs.Error(err)
			//更新任务失败 可能被杀死了 也可能是数据库错误 过5m再调用一下
			TaskRuner.Push(its, time.Now().Unix()+5)
			return
		}
		//调用去运行
		its.Calltorun(taskinfo)

	} else {

		//ready可能需要更新为pending状态，同时calltime需要更新为0
		_, err := its.ReUpdatePennding(relyinfo)
		if err != nil {
			//todo:
		}
		//依赖存在可以重试
		//增加一下运行时间 5 秒，重新去排队
		TaskRuner.Push(its, time.Now().Unix()+5)
	}

}

//rpc接口创建任务
func (p *PManage) RpcCreateTask(str string) (int64, error) {
	var rs Task
	err := json.Unmarshal([]byte(str), &rs)

	if err != nil {
		return 0, err
	}
	return p.CreateTask(rs)
}

//rpc接口更新任务
func (p *PManage) RpcUpdateTask(str string, fileds string) (bool, error) {
	var rs Task
	err := json.Unmarshal([]byte(str), &rs)

	if err != nil {
		return false, errors.New("string join params error")
	}
	m := []string{}
	if fileds != "" {
		m = strings.Split(fileds, ",")
	}
	return p.UpdateTask(rs, m)
}

//rpc接口杀死实例
func (p *PManage) RpcKillInstance(id int64) error {

	return p.killinstance(id)
}

//rpc接口标记任务成功
func (p *PManage) RpcScuccessTask(id int64) error {

	its := Instance{}
	thisit := its.GetId(id)
	if thisit == nil {
		return errors.New("instance not exsit")
	}
	_, err := thisit.Scuccess()
	return err

}

//rpc接口标记任务失败
func (p *PManage) RpcFailTask(id int64) error {
	its := Instance{}
	thisit := its.GetId(id)
	if thisit == nil {
		return errors.New("instance not exsit")
	}
	_, err := thisit.Fail()
	return err
}

//rpc接口标记任务杀死
func (p *PManage) RpcKilledTask(id int64) error {
	its := Instance{}
	thisit := its.GetId(id)
	if thisit == nil {
		return errors.New("instance not exsit")
	}
	_, err := thisit.Killed()
	return err
}

//rpc接口创建任务实例
func (p *PManage) RpcCreateTaskInstance(taskid int64, stime uint, edtime uint) (int64, error) {
	var rs Task
	thistask := rs.GetTaskbyid(taskid)
	if thistask == nil {
		return 0, errors.New("Task not exsit")
	}
	if thistask.Crontab == "" {
		return 0, errors.New("Task have no runing config")
	}

	var c int64
	thistask.UpdateNextRuntime(int64(stime))
	for thistask.Nextrun < edtime {
		go p.PushRuner(*thistask)
		thistask.UpdateNextRuntime(int64(thistask.Nextrun))
		c = c + 1
	}
	return c, nil

}
