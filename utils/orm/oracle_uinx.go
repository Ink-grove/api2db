//go:build !windows && oracle

package orm

import (
	"github.com/dzwvip/oracle"
	_ "github.com/godror/godror"
	"gorm.io/gorm"
)

// gormOracle 初始化oracle数据库
// 如果需要Oracle库 放开import里的注释 把下方 mysql.Config 改为 oracle.Config ;  mysql.New 改为 oracle.New
func gormOracle(c *config.DbConfig, m *gorm.Config) *gorm.DB {
	oracleConfig := oracle.Config{
		DSN:               c.DSN, // DSN data source name
		DefaultStringSize: 191,   // string 类型字段的默认长度
	}
	if db, err := gorm.Open(oracle.New(oracleConfig), m); err != nil {
		panic(err)
	} else {
		return db
	}
}
