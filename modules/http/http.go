package http

import (
	"crypto/tls"
	"encoding/json"
	"github.com/gogf/gf/util/gconv"
	"http2db/const"
	"http2db/models"
	"http2db/utils"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type Request interface {
	NewRequest(method, path string, body map[string]interface{}) *http.Request
}

// Caller 调用请求和解析请求
// 上层模块可以直接获取封装后的请求结果
type Caller struct {
	httpCall

	totalPage  float64
	inputParam *models.InputParam
	Client     *Client
	Request    Request
}

func NewCaller(param string, commonCall *models.CommonCall) (*Caller, error) {
	var builder Caller

	err := json.Unmarshal([]byte(param), &builder.httpCall.config)
	if err != nil {
		return nil, err
	}

	builder.httpCall.CommonCall = commonCall
	builder.Request = CallbackEntry(builder.CallMethod, builder.RequestInfo.Get())
	if builder.Request == nil {
		return nil, _const.ERR_CORE_REQUEST_NIL
	}
	builder.Client = NewHttpClient()
	if empty := builder.httpCall.isEmpty(); empty {
		return nil, _const.ERR_CORE_PARAM_EMPTY
	}

	builder.inputParam = &models.InputParam{
		ApiParam: &models.ApiParam{
			HttpMethod: builder.HttpMethod,
			Path:       builder.Path,
			Body:       builder.Body,
		},
	}

	return &builder, nil
}

func (c *Caller) Call(pageNum float64) (*models.DataReturn, *models.InputParam, error) {
	p := c.contentPacking(pageNum)

	req := c.Request.NewRequest(p.GetApiParam().GetHttpMethod(), p.GetApiParam().GetPath(), p.GetApiParam().GetBody())
	resp, err := c.Client.Call(req)
	if err != nil {
		return nil, p, err
	}
	defer resp.Body.Close()

	var tempData map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&tempData); err != nil {
		return nil, p, err
	}

	if c.totalPage < 1 {
		c.SetPageCount(tempData)
	}

	return &models.DataReturn{Data: utils.GetDataFromMapData(reflect.ValueOf(tempData), c.GetDataFiledSlice(), 0)}, nil, nil
}

func (c *Caller) GetPageCount() float64 {
	if c.SrcConfig.GetPageCount() > 0 {
		return c.SrcConfig.GetPageCount()
	}
	return c.totalPage
}

func (c *Caller) SetPageCount(data map[string]interface{}) {
	filedSlice := strings.Split(c.GetPageCountFiled(), ".")
	value := reflect.ValueOf(data)
	c.totalPage = utils.GetTotalDataFromMapData(value, filedSlice, 0)
}

func (c *Caller) contentPacking(pageNum float64) *models.InputParam {
	param := c.inputParam.Copy()

	if param.PageSize <= 0 {
		if c.SrcConfig.GetPageSize() > 0 {
			param.PageSize = c.SrcConfig.GetPageSize()
		} else {
			param.PageSize = 100
		}
	}

	switch c.GetPageFiledPosition().Type {
	case PageFiledBody:
		param.ApiParam.Body = c.apiSyncSetDataToBody(&param, pageNum)
	case PageFiledPath:
		param.ApiParam.Path = c.apiSyncSetDataToPath(&param, pageNum)
	default:
		if c.HttpMethod == http.MethodPost {
			param.ApiParam.Body = c.apiSyncSetDataToBody(&param, pageNum)
		} else if c.HttpMethod == http.MethodGet {
			param.ApiParam.Path = c.apiSyncSetDataToPath(&param, pageNum)
		}
	}

	return &param
}

func (c *Caller) apiSyncSetDataToBody(param *models.InputParam, pageNum float64) map[string]interface{} {
	param.ApiParam.Body = utils.SetFieldToMapData(
		reflect.ValueOf(param.GetApiParam().GetBody()),
		reflect.ValueOf(pageNum),
		c.GetPageFiledPosition().GetNumSlice(), 0)

	return utils.SetFieldToMapData(
		reflect.ValueOf(param.GetApiParam().GetBody()),
		reflect.ValueOf(param.PageSize),
		c.GetPageFiledPosition().GetCountSlice(), 0)
}

func (c *Caller) apiSyncSetDataToPath(param *models.InputParam, pageNum float64) string {
	path, _ := url.Parse(param.GetApiParam().GetPath())
	params := url.Values{}

	if c.GetPageFiledPosition().GetNumName() != "" {
		params.Add(c.GetPageFiledPosition().GetNumName(), gconv.String(pageNum))
	}

	if c.GetPageFiledPosition().GetCountName() != "" {
		params.Add(c.GetPageFiledPosition().GetCountName(), gconv.String(param.PageSize))
	}

	if path.RawQuery == "" {
		path.RawQuery = params.Encode()
	} else {
		path.RawQuery += "&" + params.Encode()
	}

	return path.String()
}

type Client struct {
	*http.Client
}

func NewHttpClient() *Client {
	dial := &net.Dialer{
		Timeout:   15 * time.Second,
		KeepAlive: 15 * time.Second,
	}

	transport := http.Transport{
		DialContext:       dial.DialContext,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
	}

	return &Client{
		&http.Client{
			Transport: &transport,
			Timeout:   60 * time.Second,
		},
	}
}

func (c *Client) Call(req *http.Request) (*http.Response, error) {
	return c.Do(req)
}
