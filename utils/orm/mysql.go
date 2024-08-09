package orm

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"http2db/config"
)

func gormMysql(c *config.DBConfig, m *gorm.Config) *gorm.DB {
	mysqlConfig := mysql.Config{
		DSN:                       c.DSN, // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), m); err != nil {
		panic(err)
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE=InnoDB")
		return db
	}
}
