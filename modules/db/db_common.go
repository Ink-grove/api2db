package db

import (
	"github.com/gogf/gf/os/glog"
	"gorm.io/gorm"
	"http2db/config"
	"http2db/utils/orm"
	"log"
	"strings"
)

type Client struct {
	*gorm.DB
}

func NewDbClient(config *config.DBConfig) *Client {
	return &Client{
		orm.NewORM(config),
	}
}

func (c *Client) Call(sqlStr string) ([]interface{}, error) {
	rows, err := c.Raw(sqlStr).Rows()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		glog.Error(err.Error() + ":" + sqlStr)
		return nil, err
	}

	count := len(columns)
	tableData := make([]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			col = strings.ToLower(col) // 是否需要统一转换为小写
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return tableData, nil
}
