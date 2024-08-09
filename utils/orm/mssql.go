package orm

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"http2db/config"
)

func gormMssql(c *config.DBConfig, m *gorm.Config) *gorm.DB {
	mssqlConfig := sqlserver.Config{
		DSN:               c.DSN, // DSN data source name
		DefaultStringSize: 191,   // string 类型字段的默认长度
	}
	if db, err := gorm.Open(sqlserver.New(mssqlConfig), m); err != nil {
		panic(err)
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE=InnoDB")
		return db
	}
}
