package db

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/os/glog"
	"http2db/const"
	"http2db/models"
	"strings"
)

// Caller 调用请求和解析请求
// 上层模块可以直接获取封装后的请求结果
type Caller struct {
	dbCall

	totalPage  float64
	Client     *Client
	inputParam *models.InputParam
}

func NewCaller(param string, commonCall *models.CommonCall) (*Caller, error) {
	var builder Caller

	builder.commonCall = commonCall

	err := json.Unmarshal([]byte(param), &builder.dbCall.callConfig)
	if err != nil {
		return nil, err
	}

	builder.Client = NewDbClient(&builder.dbCall.callConfig.DbConfig)
	if builder.Client == nil {
		return nil, _const.ERR_CORE_DB_ERROR
	}

	if empty := builder.dbCall.isEmpty(); empty {
		return nil, _const.ERR_CORE_PARAM_EMPTY
	}

	builder.SetPageCount()

	builder.inputParam = &models.InputParam{
		//ApiParam: &models.ApiParam{
		//	HttpMethod: builder.HttpMethod,
		//	Path:       builder.Path,
		//	Body:       builder.Body,
		//},
		PageSize: commonCall.SrcConfig.PageSize,
	}

	return &builder, nil
}

func (c *Caller) Call(pageNum float64) (*models.DataReturn, *models.InputParam, error) {
	p := c.contentPacking(pageNum)

	resp, err := c.Client.Call(p.GetDbParam().GetSql())
	if err != nil {
		return nil, p, err
	}

	return &models.DataReturn{Data: resp}, p, nil
}

func (c *Caller) contentPacking(pageNum float64) *models.InputParam {
	return &models.InputParam{
		DbParam: &models.DbParam{
			Sql: fmt.Sprintf("%s limit %d,%d", c.callConfig.SqlConfig.Sql, int(pageNum)*c.inputParam.PageSize, c.inputParam.PageSize),
		},
	}
}

func (c *Caller) GetPageCount() float64 {
	if c.commonCall.SrcConfig.GetPageCount() > 0 {
		return c.commonCall.SrcConfig.GetPageCount()
	}
	return c.totalPage
}

func (c *Caller) SetPageCount() {
	sql := c.callConfig.SqlConfig.Sql
	count := int64(0)
	fromPos := strings.LastIndex(sql, "from")

	err := c.Client.DB.Table(sql[fromPos+4:]).Count(&count).Error
	if err != nil {
		glog.Error(err)
	}

	c.totalPage = float64(count)
}
