package http

import (
	"http2db/models"
	"http2db/utils"
	"strings"
)

type httpCall struct {
	*models.CommonCall
	config
}

type config struct {
	CallMethod string                 `json:"call_method"`
	HttpMethod string                 `json:"http_method"` // http调用类型
	Path       string                 `json:"path"`
	Body       map[string]interface{} `json:"body"`

	PageCountFiled    string             `json:"page_count_filed"`    // 标识总页数的字段名，如果 PageCount
	PageFiledPosition *PageFiledPosition `json:"page_filed_position"` // 页字段位置信息，
	RequestInfo       RequestInfo        `json:"request_info"`        // request请求额外信息，账号密码信息等其它

	DataFiled      string   `json:"data_filed"`
	DataFiledSlice []string `json:"data_filed_slice"`
}

func (h *httpCall) isEmpty() bool {
	return utils.IsEmpty(h.TargetConfig.TableName, h.Filed.UniqueFieldName, h.Body, h.HttpMethod, h.Path, h.Filed.Md5FiledList)
}

func (c *config) GetPageCountFiled() string {
	return c.PageCountFiled
}
func (c *config) GetPageFiledPosition() *PageFiledPosition {
	return c.PageFiledPosition
}

func (c *config) GetDataFiledSlice() []string {
	if len(c.DataFiledSlice) < 1 {
		c.DataFiledSlice = strings.Split(c.DataFiled, ".")
	}

	return c.DataFiledSlice
}

type PageFiledPosition struct {
	Type          PageFiledType `json:"type"` // 类型，0.【智能模式-默认】Post请求放body内，get请求放path内  1.【body内】放在body体内  2.【Path内】 放在Path路径内
	NumPosition   string        `json:"num_position"`
	CountPosition string        `json:"count_position"`

	NumSlice   []string
	CountSlice []string
}

func (p *PageFiledPosition) GetType() PageFiledType {
	return p.Type
}

func (p *PageFiledPosition) GetNumName() string {
	p.initNumSlice()
	return p.NumSlice[len(p.NumSlice)-1]
}

func (p *PageFiledPosition) GetNumSlice() []string {
	p.initNumSlice()
	return p.NumSlice
}

func (p *PageFiledPosition) GetCountName() string {
	p.initCountSlice()
	return p.CountSlice[len(p.CountSlice)-1]
}

func (p *PageFiledPosition) GetCountSlice() []string {
	p.initCountSlice()
	return p.CountSlice
}

func (p *PageFiledPosition) initNumSlice() {
	if len(p.NumSlice) < 1 && p.NumPosition != "" {
		p.NumSlice = strings.Split(p.NumPosition, ".")
	}
}

func (p *PageFiledPosition) initCountSlice() {
	if len(p.CountSlice) < 1 && p.CountPosition != "" {
		p.CountSlice = strings.Split(p.CountPosition, ".")
	}
}

type PageFiledType int

const (
	PageFiledNormal = PageFiledType(0) // 【智能模式-默认】Post请求放body内，get请求放path内
	PageFiledBody   = PageFiledType(1) // 【body内】放在body体内
	PageFiledPath   = PageFiledType(2) // 【Path内】 放在Path路径内
)

type RequestInfo []byte

func (r *RequestInfo) Get() []byte {
	return *r
}

func (r *RequestInfo) UnmarshalJSON(data []byte) error {
	*r = data
	return nil
}

func (r *RequestInfo) MarshalJSON() ([]byte, error) {
	return nil, nil
}
