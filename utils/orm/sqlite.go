package orm

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"http2db/config"
)

func gormSqlite(c *config.DBConfig, m *gorm.Config) *gorm.DB {
	if db, err := gorm.Open(sqlite.Open(c.DSN), m); err != nil {
		panic(err)
	} else {
		return db
	}
}
