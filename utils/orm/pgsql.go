package orm

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"http2db/config"
)

func gormPgSql(c *config.DBConfig, m *gorm.Config) *gorm.DB {
	pgsqlConfig := postgres.Config{
		DSN:                  c.DSN, // DSN data source name
		PreferSimpleProtocol: false,
	}
	if db, err := gorm.Open(postgres.New(pgsqlConfig), m); err != nil {
		panic(err)
	} else {
		return db
	}
}
