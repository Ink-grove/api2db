package modules

import (
	"fmt"
	"github.com/golang/glog"
	"gorm.io/gorm"
	"http2db/models"
)

// 对应tool.yaml选项中 QuickFilteringMode 的 false
// 说明：即从数据库中每次进行实时查询

type Realtime struct {
	selectSql string

	db              *gorm.DB
	table           string
	pageSize        int
	md5FiledList    string
	uniqueFieldName string
}

func NewRealtimeFilter(table, md5FiledList, uniqueFieldName string, pageSize int, db *gorm.DB) *Realtime {
	return &Realtime{
		db:              db,
		table:           table,
		pageSize:        pageSize,
		md5FiledList:    md5FiledList,
		uniqueFieldName: uniqueFieldName,
	}
}

func (r *Realtime) Filter(uniqueValue string) (string, bool) {
	if r.selectSql == "" {
		r.selectSql = fmt.Sprintf("select md5(concat_ws('',%s)) as %s,%s as %s from %s ",
			r.md5FiledList, models.DefaultDbMd5FiledList[0],
			r.uniqueFieldName, models.DefaultDbMd5FiledList[1],
			r.table,
		)
	}

	var dbData models.DbMd5Data
	var limit = fmt.Sprintf(" where %s='%v' ", r.uniqueFieldName, uniqueValue)
	err := r.db.Raw(r.selectSql + limit).Scan(&dbData).Error
	if err != nil {
		glog.Error("GetFirstStructsFromSQL failed ! ", err.Error())
	}

	return dbData.Md5, dbData.Md5 != ""
}
