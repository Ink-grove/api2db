package modules

import (
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"github.com/golang/glog"
	"gorm.io/gorm"
	"http2db/models"
)

// 对应tool.yaml选项中 QuickFilteringMode 的true
// 说明：预先加载数据库中的md5值进行内存数据比对

type Memory struct {
	db              *gorm.DB
	dbMd5DataMap    map[string]string
	table           string
	pageSize        int
	md5FiledList    string
	uniqueFieldName string
}

func NewMemoryFilter(table, md5FiledList, uniqueFieldName string, pageSize int, db *gorm.DB) *Memory {
	m := &Memory{
		db:              db,
		table:           table,
		pageSize:        pageSize,
		md5FiledList:    md5FiledList,
		uniqueFieldName: uniqueFieldName,
	}

	m.initDbMd5Data()

	return m
}

func (m *Memory) Filter(uniqueValue string) (string, bool) {
	v, ok := m.dbMd5DataMap[uniqueValue]
	return v, ok
}

func (m *Memory) initDbMd5Data() {
	if m.dbMd5DataMap != nil {
		return
	}

	selectSql := fmt.Sprintf("select md5(concat_ws('',%s)) as %s,%s as %s from %s ",
		m.md5FiledList, models.DefaultDbMd5FiledList[0],
		m.uniqueFieldName, models.DefaultDbMd5FiledList[1],
		m.table,
	)

	var total int64 = 0

	m.db.Table(m.table).Count(&total)
	if total < 1 {
		return
	}

	m.dbMd5DataMap = make(map[string]string, total)

	for i := 0; i < (int(total)/m.pageSize)+1; i++ {
		var dbDataList []models.DbMd5Data
		var limit = fmt.Sprintf(" limit %v,%v", i*m.pageSize, m.pageSize)
		err := m.db.Raw(selectSql + limit).Scan(&dbDataList).Error
		if err != nil {
			glog.Error("GetTableFields failed ! ", err.Error())
		}
		m.initMd5Data(dbDataList)
	}
}

func (m *Memory) initMd5Data(data []models.DbMd5Data) {
	if m.dbMd5DataMap == nil {
		m.dbMd5DataMap = make(map[string]string)
	}

	for i := 0; i < len(data); i++ {
		m.dbMd5DataMap[gconv.String(data[i].UniqueValue)] = data[i].Md5
	}
}
