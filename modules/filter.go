package modules

import "gorm.io/gorm"

// 对应tool.yaml选项中 QuickFilteringMode

type DataFilter interface {
	Filter(uniqueValue string) (string, bool)
}

func ChooseFilter(mode bool, table, md5FiledList, uniqueFieldName string, pageSize int, db *gorm.DB) DataFilter {
	if mode {
		return NewMemoryFilter(table, md5FiledList, uniqueFieldName, pageSize, db)
	} else {
		return NewRealtimeFilter(table, md5FiledList, uniqueFieldName, pageSize, db)
	}
}
