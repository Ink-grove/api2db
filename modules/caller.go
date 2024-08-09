package modules

import (
	"errors"
	"http2db/models"
	"http2db/modules/db"
	"http2db/modules/http"
)

type CallManager interface {
	Call(pageNum float64) (returnData *models.DataReturn, record *models.InputParam, err error)

	GetPageCount() float64
}

func (t *Task) NewCallerHandle(commonCall *models.CommonCall) (CallManager, error) {
	switch t.Param.EnterParam.SyncType {
	case SyncByApi:
		return DSUEntry_Http(t.Param.EnterParam.Param, commonCall)
	case SyncByDb:
		return DSUEntry_Db(t.Param.EnterParam.Param, commonCall)
	default:
		return nil, errors.New("not found the callManager")
	}
}

func DSUEntry_Http(param string, commonCall *models.CommonCall) (CallManager, error) {
	return http.NewCaller(param, commonCall)
}

func DSUEntry_Db(param string, commonCall *models.CommonCall) (CallManager, error) {
	return db.NewCaller(param, commonCall)
}
