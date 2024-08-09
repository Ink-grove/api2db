package ctl

import (
	"context"
	"github.com/gogf/gf/os/glog"
	"github.com/robfig/cron/v3"
	"github.com/sourcegraph/conc/pool"
	"gorm.io/gorm"
	"http2db/config"
	"http2db/modules"
	"http2db/utils/orm"
	"log"
	"sync"
	"time"
)

// 数据同步器 （Data Synchronizer Unit）

type Controller struct {
	cfg *config.Config

	cron *cron.Cron
	pool *pool.Pool
	db   *gorm.DB

	dataAnalyzeSemaphore chan struct{}

	mu    sync.Mutex
	Tasks []*modules.Task
}

func (c *Controller) Init() {
	c.initConfig()
	c.initDb()
	c.initCron()
	c.initPool()
	c.initSemaphore()
}

func (c *Controller) Run() {
	var tasksInfo []modules.Param

	err := c.db.Table("hcd_dsu").Where("status = ?", 1).Find(&tasksInfo).Error
	if err != nil {
		log.Println("no tasks to run:", err.Error())
	}

	c.runAllDSU(tasksInfo)
}

func (c *Controller) AllCronJob() []cron.Entry {
	return c.cron.Entries()
}

func (c *Controller) StopAllTask() {
	for i := 0; i < len(c.Tasks); i++ {
		if c.Tasks[i].TaskId > 0 {
			c.stopDSU(c.Tasks[i])
		}
	}
}

func (c *Controller) StopSingerTask(taskId ...int) error {
	if len(taskId) >= 1 {
		t := c.getTaskByTaskId(taskId[0])
		if t == nil {
			return ERR_NOT_FOUND_TASK
		}
		c.stopDSU(t)
	}
	return nil
}

func (c *Controller) Pause(taskId ...int) error {
	if len(taskId) >= 1 {
		t := c.getTaskByTaskId(taskId[0])
		if t == nil {
			return ERR_NOT_FOUND_TASK
		}
		t.Pause()
	}
	return nil
}

func (c *Controller) KeepOn(taskId ...int) error {
	if len(taskId) >= 1 {
		t := c.getTaskByTaskId(taskId[0])
		if t == nil {
			return ERR_NOT_FOUND_TASK
		}
		t.KeepOn()
	}
	return nil
}

func (c *Controller) GetTaskStatus(taskId ...int) (float64, error) {
	if len(taskId) >= 1 {
		t := c.getTaskByTaskId(taskId[0])
		if t == nil {
			return 0, ERR_NOT_FOUND_TASK
		}
		return t.Schedule(), nil
	}
	return 0, ERR_NOT_FOUND_TASK
}

func (c *Controller) runAllDSU(paramsInfo []modules.Param) {
	glog.Info("dsu is running")

	for i := 0; i < len(paramsInfo); i++ {

		ctx, cancel := context.WithCancel(context.Background())

		if t, err := paramsInfo[i].NewTask(
			c.pool,
			c.dataAnalyzeSemaphore, ctx, cancel,
			c.cfg.DataAnalyze.MapperMaxGoroutines,
			c.cfg.SpecialChar,
			c.cfg.QuickFilteringMode,
			c.db,
		); err != nil {
			glog.Error(err.Error())
		} else {
			c.addTask(t)
		}
	}

	c.cron.Start()
}

// RunDSUAlways todo  need to del when stop ,can't let it run always
func (c *Controller) runDSUAlways(task *modules.Task) {
	task.RunAlways()
}

func (c *Controller) runDSUOnTime(task *modules.Task) {
	c.addCronJob(task, task.Param.EnterParam.Spec)
}

func (c *Controller) runDsuAndOnTime(task *modules.Task) {
	go task.Run()
	c.addCronJob(task, task.Param.EnterParam.Spec)
}

func (c *Controller) runDSUOnce(task *modules.Task) {
	go task.Run()
}

func (c *Controller) stopDSU(task *modules.Task) {
	c.cron.Stop()

	c.cron.Remove(task.CronId) // 移除cron定时任务
	task.StopTask()            // 停止当前运行中的cron任务
}

func (c *Controller) UpdateDSU() {
	// 停止之前的任务

	// 删除之前的任务

	// 将新任务加进去

	// 执行新任务
}

func (c *Controller) initCron() {
	local, _ := time.LoadLocation("Asia/Shanghai")
	c.cron = cron.New(cron.WithLocation(local), cron.WithSeconds()) // 设置时区并且精度按秒。
}

func (c *Controller) initPool() {
	c.pool = pool.New().WithMaxGoroutines(c.cfg.MaxPool)
}

func (c *Controller) initSemaphore() {
	c.dataAnalyzeSemaphore = make(chan struct{}, c.cfg.DataAnalyze.MaxSemaphore)
}

func (c *Controller) initDb() {
	db := orm.NewORM(c.cfg.DBConfig)
	if db == nil {
		panic("link db fail")
	}

	c.db = db
}

func (c *Controller) initConfig() {
	c.cfg = config.Global()
}

func (c *Controller) addTask(task *modules.Task) {
	if task == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Tasks = append(c.Tasks, task)
	c.runTask(task)
}

func (c *Controller) runTask(task *modules.Task) {
	switch task.Param.EnterParam.DsuType {
	case modules.DsuRunAlways:
		c.runDSUAlways(task)
	case modules.DsuRunOnTime:
		c.runDSUOnTime(task)
	case modules.DsuRunAndOnTime:
		c.runDsuAndOnTime(task)
	case modules.DsuRunOnce:
		c.runDSUOnce(task)
	default:

	}
}

func (c *Controller) addCronJob(task *modules.Task, spec string) {
	var err error
	task.CronId, err = c.cron.AddJob(spec, task)
	if err != nil {
		glog.Error("addCronJob failed:", err.Error())
	}
}

func (c *Controller) getTaskByTaskId(taskId int) *modules.Task {
	for i := 0; i < len(c.Tasks); i++ {
		if c.Tasks[i].TaskId == taskId {
			return c.Tasks[i]
		}
	}
	return nil
}
