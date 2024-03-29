package taskgroup

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sjlleo/traceSysBackend/email"
	"github.com/sjlleo/traceSysBackend/models"
)

var ActiveSchedulers map[uint]Task

type Task struct {
	TaskDetail      models.Tasks
	IntervalSeconds int
	Scheduler       *gocron.Scheduler
	ResultRWLock    *sync.RWMutex // 读写互斥锁
}

func init() {
	ActiveSchedulers = make(map[uint]Task)
}

func NewTask() *Task {
	return &Task{
		Scheduler: gocron.NewScheduler(time.UTC),
	}
}

func StartTaskCycle() {
	t := NewTask()
	log.Println("TaskGroup - StartTaskCycle")
	t.Scheduler.Every("10m").Do(t.getTaskList)
	t.Scheduler.StartBlocking()
}

func (t *Task) getTaskList() {
	// 向数据库获取任务清单
	list, err := models.GetAllTasks()
	if err != nil {
		log.Println(err)
		return
	}

	go t.UpdateScheduler(list)
	go t.CleanUpScheduler(list)
}

func (t *Task) DoTask() bool {
	var model, msg string
	needSend := models.CheckExceed(&t.TaskDetail)
	if needSend {
		switch t.TaskDetail.TraceType {
		case 1:
			// 获取模板
			temp, err := models.GetTemplate(models.RTTExceed)
			if err != nil {
				log.Println(err)
				return true
			}
			model = temp.Model
		case 2:
			// 获取模板
			temp, err := models.GetTemplate(models.PacketLossExceed)
			if err != nil {
				log.Println(err)
				return true
			}
			model = temp.Model

		}
		tag, err := models.FindTargetIPByID(t.TaskDetail.TargetID)
		if err != nil {
			log.Println(err)
			return true
		}
		msg = strings.Replace(model, "{ip}", tag.TargetIP, -1)
		// 查找用户的信息
		u, _ :=models.FindUserByID(t.TaskDetail.CreatedUserID)
		if t.TaskDetail.CallMethod == 1 {
			email.SendMsg(msg, u.Email)
		}
		return true
	}
	return false
}

func (t *Task) UpdateScheduler(list []models.Tasks) {
	log.Println("traceService - UpdateScheduler")
	for _, t := range list {
		if _, ok := ActiveSchedulers[t.ID]; !ok {
			// 未开启的创建新的 Schedulers
			ActiveSchedulers[t.ID] = Task{
				TaskDetail:   t,
				Scheduler:    gocron.NewScheduler(time.UTC),
				ResultRWLock: &sync.RWMutex{},
			}
			// 获得 Map 对应 Key 的地址
			s := ActiveSchedulers[t.ID]
			s.Scheduler.Every(60).Minutes().Do(s.DoTask)
			s.Scheduler.StartAsync()
		}
	}
}

func (t *Task) CleanUpScheduler(list []models.Tasks) {
	log.Println("traceService - CleanUpScheduler")

	for taskID, activeTask := range ActiveSchedulers {
		var taskShouldDelete bool = true
		for _, pendingTask := range list {
			if pendingTask.ID == taskID {
				taskShouldDelete = false
				break
			}
		}

		if taskShouldDelete {
			activeTask.Scheduler.Stop()
			delete(ActiveSchedulers, taskID)
		}
	}
}
