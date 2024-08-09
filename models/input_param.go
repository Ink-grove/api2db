package models

// InputParam api sync请求入参
type InputParam struct {
	ApiParam *ApiParam
	DbParam  *DbParam
	PageSize int
}

func (a *InputParam) Copy() InputParam {
	return InputParam{
		ApiParam: a.ApiParam,
		DbParam:  a.DbParam,
		PageSize: a.PageSize,
	}
}

func (a *InputParam) GetApiParam() *ApiParam {
	return a.ApiParam
}

func (a *InputParam) GetDbParam() *DbParam {
	return a.DbParam
}

func (a *InputParam) GetPageSize() int {
	return a.PageSize
}

type ApiParam struct {
	HttpMethod string
	Path       string
	Body       map[string]interface{}
}

func (a *ApiParam) GetHttpMethod() string {
	return a.HttpMethod
}

func (a *ApiParam) GetPath() string {
	return a.Path
}

func (a *ApiParam) GetBody() map[string]interface{} {
	return a.Body
}

type DbParam struct {
	Sql string
}

func (d *DbParam) GetSql() string {
	return d.Sql
}
