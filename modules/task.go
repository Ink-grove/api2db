package modules

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"github.com/robfig/cron/v3"
	"github.com/sourcegraph/conc/iter"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/atomic"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"http2db/const"
	"http2db/models"
	"http2db/utils"
	"reflect"
	"strings"
	"sync"
)

// Task 具体任务信息
type Task struct {
	TaskId int

	db              *gorm.DB
	ctx             context.Context
	cancel          context.CancelFunc
	finalDataDeal   finalDataDealFunc // 最终数据处理函数
	filter          DataFilter
	pauseSignalChan chan struct{}

	CompleteCount atomic.Float64 // 当前完成页数
	PageCount     float64        // 总页数

	uniqueFieldName       reflect.Value
	compareFiledListMap   map[string]reflect.Value
	saveFiledListWithAlia []string
	md5FiledListWithAlia  []string

	srcMd5DataMap sync.Map

	dataAnalyzeSemaphore chan struct{}
	mapper               iter.Mapper[interface{}, map[string]interface{}]
	state                State // 任务状态
	stateMux             sync.RWMutex
	pool                 *pool.Pool
	CronId               cron.EntryID `json:"cron_id"`
	Param                *TaskParam   `json:"param"`
}

type finalDataDealFunc func(object interface{}) interface{}

func (p *Param) NewTask(pool *pool.Pool, c chan struct{}, ctx context.Context, cancel context.CancelFunc, max int, specialChar, filterMode bool, db *gorm.DB) (*Task, error) {
	var callInfo *models.CommonCall
	var f finalDataDealFunc

	err := json.Unmarshal([]byte(p.InitParam), &callInfo)
	if err != nil {
		return nil, _const.ERR_CORE_UNMARSHAL_ERROR
	}

	if specialChar {
		f = utils.SpecialDataDeal
	} else {
		f = utils.FormalDataDeal
	}

	table := callInfo.GetTargetConfig().GetTableName()
	pageSize := callInfo.GetTargetConfig().GetPageSize()
	md5FiledList := callInfo.GetFiled().GetMd5FiledList()
	uniqueFieldName := callInfo.GetFiled().GetUniqueFieldName()

	return &Task{
		TaskId:               p.Id,
		db:                   db,
		finalDataDeal:        f,
		filter:               ChooseFilter(filterMode, table, md5FiledList, uniqueFieldName, pageSize, db),
		ctx:                  ctx,
		cancel:               cancel,
		pauseSignalChan:      make(chan struct{}, 1),
		pool:                 pool,
		mapper:               iter.Mapper[interface{}, map[string]interface{}]{MaxGoroutines: max},
		dataAnalyzeSemaphore: c,
		Param: &TaskParam{
			CallInfo:   callInfo,
			EnterParam: p,
		},
	}, nil
}

// Status 任务状态
func (t *Task) Status() *State {
	cIdx := t.CompleteCount.Load()
	if t.PageCount <= cIdx {
		t.state = CompletedState
	}
	return t.state.GetState()
}

// Schedule 进度情况
func (t *Task) Schedule() (score float64) {
	cIdx := t.CompleteCount.Load()
	if t.PageCount > 1 && cIdx > 1 {
		score = cIdx / t.PageCount
		if score >= 1 {
			return 1
		}
		return
	}
	return 0
}

// Run 定时运行
func (t *Task) Run() {
	t.execute()
}

// RunAlways 永久运行
func (t *Task) RunAlways() {
	go func() {
		for {
			t.execute()
		}
	}()
}

func (t *Task) StopTask() {
	t.cancel()
}

func (t *Task) Pause() {
	if t.loadState() == RunningState {
		t.pauseSignalChan <- struct{}{}
		t.changeState(PausedState)
		glog.Info(fmt.Sprintf("task id :%v is pause", t.TaskId))
	}
}

func (t *Task) KeepOn() {
	if t.loadState() == PausedState {
		<-t.pauseSignalChan
		t.changeState(RunningState)
		glog.Info(fmt.Sprintf("task id :%v is keep on", t.TaskId))
	}
}

func (t *Task) execute() {
	if !(t.loadState() == StartState || t.loadState() == CompletedState) {
		return
	}

	t.changeState(RunningState)
	t.init()
	i := t.getFirstPageNum()

	handler, err := t.NewCallerHandle(t.Param.CallInfo)
	if err != nil {
		glog.Error(err)
		return
	}

	returnData, record, err := handler.Call(i)
	if err != nil {
		t.recordFailedTasks(record, err)
		return
	}

	t.PageCount = handler.GetPageCount() // 获取到总页数

	err = t.dataGovernance(returnData)
	if err != nil {
		t.recordFailedTasks(record, err)
		return
	}

	// 执行每页遍历操作
	for i++; i <= t.PageCount; i++ {
		newReturnData, record, err := handler.Call(i)
		if err != nil {
			t.recordFailedTasks(record, err)
		}

		select {
		case <-t.ctx.Done():
			t.changeState(ClosedState)
			glog.Info(fmt.Sprintf("receive done , kill the execute , task id :%v", t.TaskId))
			return
		default:
			t.pauseSignalChan <- struct{}{}
			<-t.pauseSignalChan
			t.pool.Go(func() {
				err := t.dataGovernance(newReturnData)
				if err != nil {
					t.recordFailedTasks(record, err)
					return
				}
			})
		}
	}

	t.changeState(CompletedState)
	glog.Info(fmt.Sprintf("task id :%v is down", t.TaskId))
}

func (t *Task) init() {
	t.srcMd5DataMap = sync.Map{}
	t.initMapper()
	t.initCompareFiledListMap()
	t.initSaveFiledListWithAlia()
	t.initUniqueFieldName()
	t.initMd5FiledListWithAlia()
}

func (t *Task) dataGovernance(data *models.DataReturn) error {
	defer t.CompleteCount.Add(1)
	if data.IsEmpty() {
		return nil
	}

	insertData := t.analyzeData(data.Data)
	return t.saveData(insertData)
}

func (t *Task) changeState(s State) {
	t.stateMux.Lock()
	defer t.stateMux.Unlock()
	t.state = s
}

func (t *Task) loadState() State {
	t.stateMux.RLock()
	defer t.stateMux.RUnlock()
	return t.state
}

func (t *Task) updateData(updateSql string) bool {
	if updateSql == "" {
		return true
	}

	err := t.db.Raw(updateSql).Error
	if err != nil {
		glog.Error("save error!", err.Error())
	}

	return false
}

func (t *Task) saveData(insertData []map[string]interface{}) error {
	if len(insertData) < 1 {
		return nil
	}

	// todo 加入批量插入+限流

	err := t.db.Table(t.Param.CallInfo.GetTargetConfig().GetTableName()).Create(insertData).Error
	if err != nil {
		return err
	}

	return nil
}

func (t *Task) analyzeData(srcData []interface{}) (insertMap []map[string]interface{}) {
	t.dataAnalyzeSemaphore <- struct{}{}
	//glog.Info("当前限制信号量中的数据量：", len(t.findNeedDataSemaphore))

	insertMap = t.mapper.Map(srcData, func(item *interface{}) (finalMapData map[string]interface{}) {
		object := reflect.ValueOf(*item)
		if object.Kind() == reflect.Map {
			uniqueField := object.MapIndex(t.uniqueFieldName)
			uniqueValue := gconv.String(uniqueField.Interface())
			finalMapData = t.getFinalDataFromObj(object)
			if _, srcExistItem := t.srcMd5DataMap.LoadOrStore(uniqueValue, ""); !srcExistItem {
				if finalMapData != nil {
					if dbMd5, dbExistItem := t.filter.Filter(uniqueValue); dbExistItem {
						md5ValueSlice := t.getMd5ValueSlice(finalMapData)
						srcMd5 := utils.Md5Encrypt(strings.Join(md5ValueSlice, ""))
						t.srcMd5DataMap.Store(uniqueValue, srcMd5)
						if srcMd5 != dbMd5 {
							saveSlice := t.getSaveSlice(finalMapData)

							t.updateData(fmt.Sprintf("update %s set %s where %s = '%s'",
								t.Param.CallInfo.GetTargetConfig().GetTableName(),
								strings.Join(saveSlice, ","),
								t.compareFiledListMap[t.Param.CallInfo.GetFiled().GetUniqueFieldName()], uniqueValue,
							))
						}
					} else {
						return finalMapData
					}
				}

			} else {
				if finalMapData != nil {
					if dbMd5, dbExistItem := t.filter.Filter(uniqueValue); dbExistItem {
						md5ValueSlice := t.getMd5ValueSlice(finalMapData)
						srcMd5 := utils.Md5Encrypt(strings.Join(md5ValueSlice, ""))
						t.srcMd5DataMap.Store(uniqueValue, srcMd5)
						if srcMd5 != dbMd5 {
							saveSlice := t.getSaveSlice(finalMapData)

							t.updateData(fmt.Sprintf("update %s set %s where %s = '%s'",
								t.Param.CallInfo.GetTargetConfig().GetTableName(),
								strings.Join(saveSlice, ","),
								t.compareFiledListMap[t.Param.CallInfo.GetFiled().GetUniqueFieldName()], uniqueValue,
							))
						}
					}
				}
			}
		}
		return nil
	})

	<-t.dataAnalyzeSemaphore
	return
}

func (t *Task) getFinalDataFromObj(object reflect.Value) map[string]interface{} {
	var (
		flag       = false
		storeValue interface{}
		tempMap    = make(map[string]interface{}, len(t.saveFiledListWithAlia))
	)

	for _, v := range t.saveFiledListWithAlia {
		filed := object.MapIndex(reflect.ValueOf(v))
		if filed.IsValid() {
			flag = true
			storeValue = nil
			if !filed.IsZero() {
				storeValue = t.finalDataDeal(filed.Interface())
			}
			if key, ok := t.compareFiledListMap[v]; ok {
				tempMap[gconv.String(key.Interface())] = storeValue
			}
		}
	}

	if flag {
		return tempMap
	}

	return nil
}

func (t *Task) getMd5ValueSlice(obj map[string]interface{}) (md5ValueSlice []string) {
	for _, v := range t.md5FiledListWithAlia {
		if key, ok := t.compareFiledListMap[v]; ok {
			value := obj[gconv.String(key.Interface())]
			if !utils.IsEmpty(value) {
				md5ValueSlice = append(md5ValueSlice, gconv.String(value))
			}
		}

	}
	return
}

func (t *Task) getSaveSlice(obj map[string]interface{}) (saveSlice []string) {
	for _, v := range t.saveFiledListWithAlia {
		if key, ok := t.compareFiledListMap[v]; ok {
			value := obj[gconv.String(key.Interface())]
			if !utils.IsEmpty(value) {
				saveSlice = append(saveSlice, fmt.Sprintf("%s = '%s'", key.Interface(), gconv.String(value)))
			} else {
				saveSlice = append(saveSlice, fmt.Sprintf("%s = null", v))
			}
		}
	}
	return
}

func (t *Task) recordFailedTasks(param *models.InputParam, err error) {
	h := &HcdDsuFailTask{
		ProjectId: t.Param.EnterParam.ProjectId,
		SyncType:  t.Param.EnterParam.SyncType,
		Name:      t.Param.EnterParam.Name,
		TaskId:    t.TaskId,
		CallInfo: map[string]interface{}{
			"client_param": t.Param.EnterParam.InitParam,
			"input_param":  param,
		},
		Discription: err.Error(),
	}

	saveErr := t.db.Table("hcd_dsu_fail_task").Save(h).Error
	if saveErr != nil {
		glog.Error("recordFailedTasks save error!", err.Error())
	}
}

func (t *Task) getFirstPageNum() float64 {
	return t.Param.CallInfo.SrcConfig.GetFirstPageNum()
}

// initCompareFiledListMap 全量的比对map，key为目标存储数据库的字段名的反射值，value为对应数据源的字段名
func (t *Task) initCompareFiledListMap() {
	if t.compareFiledListMap != nil {
		return
	}

	t.compareFiledListMap = make(map[string]reflect.Value, 10)

	filedSlice := t.Param.CallInfo.GetFiled().GetSaveFiledSlice()
	if filedSlice == nil {
		filedSlice = t.GetDbTableFields(t.Param.CallInfo.GetTargetConfig().GetTableName())
	}

	for _, v := range filedSlice {
		t.compareFiledListMap[v] = reflect.ValueOf(v)
	}

	filedAlias := t.Param.CallInfo.GetFiled().GetFiledAlias()

	for k, v := range filedAlias {
		t.compareFiledListMap[v] = reflect.ValueOf(k)
	}

}

func (t *Task) initSaveFiledListWithAlia() {
	if t.saveFiledListWithAlia != nil || len(t.saveFiledListWithAlia) > 1 {
		return
	}

	var tempMap = make(map[string]struct{})
	var filedSlice []string
	filed := t.Param.CallInfo.GetFiled().GetSaveFiledList()
	if filed == "" {
		filedSlice = t.GetDbTableFields(t.Param.CallInfo.GetTargetConfig().GetTableName())
	} else {
		filedSlice = t.Param.CallInfo.GetFiled().GetSaveFiledSlice()
	}

	for _, v := range filedSlice {
		if v != "" {
			tempMap[v] = struct{}{}
		}
	}

	filedAlias := t.Param.CallInfo.GetFiled().GetFiledAlias()

	for k, v := range filedAlias {
		if _, ok := tempMap[k]; ok {
			delete(tempMap, k)
			tempMap[v] = struct{}{}
		}
	}

	for k, _ := range tempMap {
		t.saveFiledListWithAlia = append(t.saveFiledListWithAlia, k)
	}

}

func (t *Task) initUniqueFieldName() {
	if t.uniqueFieldName.IsValid() {
		return
	}
	uniqueFieldName := t.Param.CallInfo.GetFiled().GetUniqueFieldName()

	filedAlias := t.Param.CallInfo.GetFiled().GetFiledAlias()

	for k, v := range filedAlias {
		if uniqueFieldName == k {
			t.uniqueFieldName = reflect.ValueOf(v)
			break
		}
	}

	if !t.uniqueFieldName.IsValid() {
		t.uniqueFieldName = reflect.ValueOf(uniqueFieldName)
	}
}

func (t *Task) initMd5FiledListWithAlia() {
	if len(t.md5FiledListWithAlia) < 1 {
		return
	}

	var tempMap = make(map[string]struct{})

	filedSlice := t.Param.CallInfo.GetFiled().GetMd5FiledSlice()
	for _, v := range filedSlice {
		if v != "" {
			tempMap[v] = struct{}{}
		}
	}

	filedAlias := t.Param.CallInfo.GetFiled().GetFiledAlias()

	for k, v := range filedAlias {
		if _, ok := tempMap[k]; ok {
			delete(tempMap, k)
			tempMap[v] = struct{}{}
		}
	}

	for k, _ := range tempMap {
		t.md5FiledListWithAlia = append(t.md5FiledListWithAlia, k)
	}

}

func (t *Task) initMapper() {
	if t.mapper.MaxGoroutines < 0 {
		t.mapper = iter.Mapper[interface{}, map[string]interface{}]{
			MaxGoroutines: 10,
		}
	}
}

func (t *Task) GetDbTableFields(table string) []string {
	tabList, err := t.db.Migrator().ColumnTypes(table)
	if err != nil {
		glog.Error("GetTableFields failed ! ", err.Error())
	}

	result := make([]string, len(tabList))
	for i := 0; i < len(tabList); i++ {
		if !utils.FindExclusionField(tabList[i].Name()) {
			result[i] = tabList[i].Name()
		}
	}

	return result
}
