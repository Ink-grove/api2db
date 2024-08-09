//go:build windows || linux

package orm

import (
	"gorm.io/gorm"
	"http2db/config"
)

// gormOracle 由于依赖库不支持windows编译，故windows编译返回nil
func gormOracle(c *config.DBConfig, m *gorm.Config) *gorm.DB {
	return nil
}
