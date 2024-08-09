package ctl

import (
	"errors"
)

var (
	ERR_NOT_FOUND_TASK = errors.New("not find the task by target taskId")
)
