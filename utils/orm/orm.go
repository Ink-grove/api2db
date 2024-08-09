package orm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"http2db/config"
	"log"
	"os"
	"time"
)

// NewORM new db and retry connection when has error.
func NewORM(c *config.DBConfig) (db *gorm.DB) {
	if c.DSN == "" || c.DbType == "" {
		return nil
	}

	db = initDb(c)

	tmp, err := time.ParseDuration(c.IdleTimeout)
	if err != nil {
		//log.Error("idleTimeout Parse error: %v", c.DSN, err)
		log.Printf("idleTimeout Parse error: %v", err)
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil
	}

	sqlDB.SetMaxIdleConns(c.Idle)
	sqlDB.SetMaxOpenConns(c.Active)
	sqlDB.SetConnMaxLifetime(tmp / time.Second)
	return
}

func initDb(c *config.DBConfig) *gorm.DB {
	m := configSetting(c)
	switch c.DbType {
	case "mysql":
		return gormMysql(c, m)
	case "pgsql":
		return gormPgSql(c, m)
	case "oracle":
		return gormOracle(c, m)
	case "mssql":
		return gormMssql(c, m)
	case "sqlite":
		return gormSqlite(c, m)
	default:
		return gormMysql(c, m)
	}
}

func configSetting(c *config.DBConfig) *gorm.Config {
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	_default := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})

	config.Logger = _default.LogMode(logger.Silent)

	return config
}
