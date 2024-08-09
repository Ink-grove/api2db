package models

type DbMd5Data struct {
	Md5         string `json:"md5"`
	UniqueValue int    `json:"unique_value"`
}

var DefaultDbMd5FiledList = []string{"md5", "unique_value"}
