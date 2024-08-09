package modules

import (
	"github.com/gogf/gf/os/gtime"
	"http2db/models"
)

// Param hcd_dsu表结构
type Param struct {
	Id          int         `json:"id"`          //
	ProjectId   int         `json:"project_id"`  // 项目
	SyncType    SyncType    `json:"sync_type"`   // 同步器类型：1.【api同步器】  2.【db同步器】
	Name        string      `json:"name"`        // 任务名称
	InitParam   string      `json:"init_param"`  // 初始化参数
	Param       string      `json:"param"`       // 执行参数
	Status      int         `json:"status"`      // 状态
	DsuType     DsuType     `json:"dsu_type"`    // 数据同步类型：1.【永久运行】  2.【按时执行】依照cron指定的参数  3.【按时执行并且启动的时候执行一次】 4.【仅启动时执行】
	DsuMode     DsuMode     `json:"dsu_mode"`    // 数据同步模式：1.【正序运行】  2.【倒序运行】 3.【中序向前运行】  4.【中序向后运行】
	Spec        string      `json:"spec"`        // cron执行参数
	Discription string      `json:"discription"` // 任务描述
	LastTime    *gtime.Time `json:"last_time"`   // 上次执行时间
	CreateTime  *gtime.Time `json:"create_time"` // 添加时间
	UpdateTime  *gtime.Time `json:"update_time"` // 更新时间
}

type TaskParam struct {
	CallInfo   *models.CommonCall
	EnterParam *Param
}

type SyncType int

const (
	SyncByApi    = SyncType(1)
	SyncByDb     = SyncType(2)
	SyncByOracle = SyncType(3)
)

type DsuType int

const (
	DsuRunAlways    = DsuType(1) // 【永久运行】
	DsuRunOnTime    = DsuType(2) // 【按时执行】
	DsuRunAndOnTime = DsuType(3) // 【按时执行并且启动的时候执行一次】
	DsuRunOnce      = DsuType(4) // 【仅启动时执行】
)

type DsuMode int

const (
	DsuOrder               = DsuMode(1) // 【正序运行】
	DsuReverseOrder        = DsuMode(2) // 【倒序运行】
	DsuMiddleOrderForward  = DsuMode(3) // 【中序向前运行】
	DsuMiddleOrderBackward = DsuMode(4) // 【中序向后运行】
)

type State int32

const (
	StartState     = State(0) // 【任务开始执行】
	RunningState   = State(1) // 【任务进行中】
	PausedState    = State(2) // 【任务已暂停】
	CompletedState = State(3) // 【任务已完成】
	ClosedState    = State(4) // 【任务已关闭】
)

func (s *State) GetState() *State {
	return s
}

// HcdDsuFailTask 数据同步器 - 失败任务表
type HcdDsuFailTask struct {
	ProjectId   int         `json:"project_id"`  // 项目
	SyncType    SyncType    `json:"sync_type"`   // 同步器类型
	Name        string      `json:"name"`        // 任务名称
	TaskId      int         `json:"task_id"`     // 任务id
	CallInfo    interface{} `json:"call_Info"`   // 调用信息，保存请求的path，body以及服务器参数等，用于任务补偿
	Discription string      `json:"discription"` // 任务错误描述
}
