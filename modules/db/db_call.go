package db

import (
	"http2db/config"
	"http2db/models"
	"http2db/utils"
)

type dbCall struct {
	commonCall *models.CommonCall
	callConfig
}

// config orm config
//type config struct {
//	Sql      string     `json:"sql"`
//	DbConfig orm.Config `json:"db_config"`
//}

type callConfig struct {
	SqlConfig SqlConfig       `json:"sql_config"`
	DbConfig  config.DBConfig `json:"db_config"`
}

type SqlConfig struct {
	Sql string `json:"sql"` // sql语句
}

func (d *dbCall) isEmpty() bool {
	return utils.IsEmpty(d.callConfig.SqlConfig.Sql, d.callConfig.DbConfig.DbType)
}
