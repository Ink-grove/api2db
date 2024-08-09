package api

import (
	"github.com/gogf/gf/net/ghttp"
	"http2db/cmd"
	"http2db/const"
	"http2db/utils"
)

// Action 公共模块接口
type Action struct {
	utils.BaseRouter
}

func NewAction() *Action {
	return &Action{}
}

// GetTaskStatus 获取Task的进度情况
func (a *Action) GetTaskStatus(r *ghttp.Request) {
	taskId := r.GetInt("task_id")
	if taskId < 1 {
		utils.Info(r, _const.ERR_NIL_TASK_ID)
	}

	if state, err := cmd.G_dsu.GetTaskStatus(taskId); err != nil {
		utils.Info(r, err.Error())
	} else {
		utils.Data(r, state)
	}
}

func (a *Action) GetAllCronJob(r *ghttp.Request) {
	result := cmd.G_dsu.AllCronJob()
	utils.Data(r, result)
}

func (a *Action) Stop(r *ghttp.Request) {
	taskId := r.GetInt("task_id")
	if taskId < 1 {
		utils.Info(r, _const.ERR_NIL_TASK_ID)
	}

	if err := cmd.G_dsu.StopSingerTask(taskId); err != nil {
		utils.Info(r, err.Error())
	} else {
		utils.InfoSuccess(r)
	}
}

func (a *Action) StopAll(r *ghttp.Request) {
	cmd.G_dsu.StopAllTask()
	utils.InfoSuccess(r)
}

func (a *Action) Pause(r *ghttp.Request) {
	taskId := r.GetInt("task_id")
	if taskId < 1 {
		utils.Info(r, _const.ERR_NIL_TASK_ID)
	}

	if err := cmd.G_dsu.Pause(taskId); err != nil {
		utils.Info(r, err.Error())
	} else {
		utils.InfoSuccess(r)
	}
}

func (a *Action) KeepOn(r *ghttp.Request) {
	taskId := r.GetInt("task_id")
	if taskId < 1 {
		utils.Info(r, _const.ERR_NIL_TASK_ID)
	}

	if err := cmd.G_dsu.KeepOn(taskId); err != nil {
		utils.Info(r, err.Error())
	} else {
		utils.InfoSuccess(r)
	}
}
