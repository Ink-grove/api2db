package _const

import "errors"

var (
	ERR_CORE_PARAM_EMPTY     = errors.New("param is empty , please check")
	ERR_CORE_UNMARSHAL_ERROR = errors.New("unmarshal error")
	ERR_CORE_REQUEST_NIL     = errors.New("request is nil , please check")

	ERR_CORE_DB_ERROR = errors.New("init db client failed")

	ERR_NIL_TASK_ID = errors.New("task id is nil")
)
