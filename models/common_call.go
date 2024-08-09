package models

import "strings"

type CommonCall struct {
	TargetConfig *TargetConfig `json:"target_config"`
	Filed        *Filed        `json:"filed"`
	SrcConfig    *SrcConfig    `json:"src_config"`
}

func (c *CommonCall) GetTargetConfig() *TargetConfig {
	return c.TargetConfig
}

func (c *CommonCall) GetFiled() *Filed {
	return c.Filed
}

func (c *CommonCall) GetSrcConfig() *SrcConfig {
	return c.SrcConfig
}

type Filed struct {
	UniqueFieldName string            `json:"unique_field_name"` // 表示唯一值的字段名
	Md5FiledList    string            `json:"md5_filed_list"`    // 需要计算md5值的数据库字段集合
	SaveFiledList   string            `json:"save_filed_list"`   // 保存数据的表字段集合
	FiledAlias      map[string]string `json:"filed_alias"`       // 字段别名,key为目标数据库表中的字段，value为对应源数据中的字段

	Md5FiledSlice  []string
	SaveFiledSlice []string
}

func (d *Filed) GetUniqueFieldName() string {
	return d.UniqueFieldName
}

func (d *Filed) GetMd5FiledList() string {
	return d.Md5FiledList
}

func (d *Filed) GetSaveFiledList() string {
	return d.SaveFiledList
}

func (d *Filed) GetFiledAlias() map[string]string {
	return d.FiledAlias
}

func (d *Filed) GetMd5FiledSlice() []string {
	if d.GetMd5FiledList() == "" {
		return nil
	}

	if d.Md5FiledSlice == nil || len(d.Md5FiledSlice) < 1 {
		if d.GetMd5FiledList() != "" {
			d.Md5FiledSlice = strings.Split(d.GetMd5FiledList(), ",")
		}
	}

	return d.Md5FiledSlice
}

func (d *Filed) GetSaveFiledSlice() []string {
	if d.GetSaveFiledList() == "" {
		return nil
	}

	if d.SaveFiledSlice == nil || len(d.SaveFiledSlice) < 1 {
		if d.GetSaveFiledList() != "" {
			d.SaveFiledSlice = strings.Split(d.GetSaveFiledList(), ",")
		}
	}
	return d.SaveFiledSlice
}

type TargetConfig struct {
	TableName string `json:"table_name"`
	PageSize  int    `json:"page_size"` // 每次查询的分页大小，主要为了查询数据的md5值
}

func (t *TargetConfig) GetTableName() string {
	return t.TableName
}

func (t *TargetConfig) GetPageSize() int {
	if t.PageSize < 1 {
		return 100000
	}
	return t.PageSize
}

type SrcConfig struct {
	PageMode     int     `json:"page_mode"`      // 页数模式，0：正常翻页模式【默认】，pageNum累加 ， 1：偏移量模式，limit + offset，根据pageNum计算offset
	FirstPageNum float64 `json:"first_page_num"` // 首次调用的pageNum
	PageCount    float64 `json:"page_count"`     // 指定遍历的总页数
	PageSize     int     `json:"page_size"`
}

func (s *SrcConfig) GetPageMode() int {
	return s.PageMode
}

func (s *SrcConfig) GetFirstPageNum() float64 {
	return s.FirstPageNum
}

func (s *SrcConfig) GetPageCount() float64 {
	return s.PageCount
}

func (s *SrcConfig) GetPageSize() int {
	return s.PageSize
}
