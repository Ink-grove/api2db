package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

func IsEmpty(values ...interface{}) bool {
	for _, value := range values {
		if value == nil {
			return true
		}

		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
			if v.Len() == 0 {
				return true
			}
		case reflect.Bool:
			if !v.Bool() {
				return true
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v.Int() == 0 {
				return true
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if v.Uint() == 0 {
				return true
			}
		case reflect.Float32, reflect.Float64:
			if v.Float() == 0 {
				return true
			}
		case reflect.Complex64, reflect.Complex128:
			if v.Complex() == 0 {
				return true
			}
		}

		if v.Kind() == reflect.String {
			if strings.ToLower(v.String()) == "null" {
				return true
			}
		}
	}

	return false
}

func DataExistInMap(value string, theMap map[reflect.Value]string) bool {
	for _, saveFiled := range theMap {
		if saveFiled == value {
			return true
		}
	}

	return false
}

func SetFieldToMapData(value, data reflect.Value, filedSlice []string, level int) map[string]interface{} {
	if t, ok := value.Interface().(map[string]interface{}); ok {
		v := reflect.ValueOf(t)
		objectField := v.MapIndex(reflect.ValueOf(filedSlice[level]))
		if len(filedSlice) == level+1 {
			v.SetMapIndex(reflect.ValueOf(filedSlice[level]), data)
		} else {
			childTerm := make(map[string]interface{})
			if objectField.IsValid() {
				childTerm = SetFieldToMapData(objectField, data, filedSlice, level+1)
				v.SetMapIndex(reflect.ValueOf(filedSlice[level]), reflect.ValueOf(childTerm))
			}
		}
		return v.Interface().(map[string]interface{})
	}
	return nil
}

func GetTotalDataFromMapData(value reflect.Value, filedSlice []string, level int) float64 {
	if value.Kind() == reflect.Map {
		allPageValue := value.MapIndex(reflect.ValueOf(filedSlice[level]))
		if allPageValue.IsValid() {
			if convertedData, ok := allPageValue.Interface().(map[string]interface{}); ok {
				return GetTotalDataFromMapData(reflect.ValueOf(convertedData), filedSlice, level+1)
			} else {
				return gconv.Float64(allPageValue.Interface())
			}
		}
	}

	return 0
}

// GetDataFromMapData value为源数据的反射值，filedSlice为取数据的切片，level表示层级
func GetDataFromMapData(value reflect.Value, filedSlice []string, level int) []interface{} {
	if value.Kind() == reflect.Map {
		allPageValue := value.MapIndex(reflect.ValueOf(filedSlice[level]))
		if allPageValue.IsValid() {
			if convertedData, ok := allPageValue.Interface().(map[string]interface{}); ok {
				return GetDataFromMapData(reflect.ValueOf(convertedData), filedSlice, level+1)
			} else if convertedData, ok := allPageValue.Interface().([]interface{}); ok {
				return convertedData
			}
		}
	}

	return nil
}

// 进行字段排除，将单引号，斜线等进行特殊处理，避免存库时出错问题产生
var dataDealReg = regexp.MustCompile(`'|\\|\s+\+0800\s+CST`)

func FormalDataDeal(values interface{}) interface{} {
	return values
}

func SpecialDataDeal(values interface{}) interface{} {
	v := gconv.String(values)
	return dataDealReg.ReplaceAllString(v, "")
}

// 进行字段排除
var fieldExclusionReg = regexp.MustCompile(``)

func FindExclusionField(s string) bool {
	s = strings.ToLower(s)
	if fieldExclusionReg.FindString(s) != "" {
		return true
	}
	return false
}

func Md5Encrypt(content string) string {
	data := []byte(content)
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func GetAppPath() string {
	dir, err := os.Executable()
	if err != nil {
		glog.Error(err)
		return ""
	}
	dir = strings.Replace(dir, "\\", "/", -1)
	filePath, _ := filepath.Split(dir)
	return filePath
}

func ToBuffer(v interface{}) *bytes.Buffer {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	// 这里要设置禁止html转义
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		glog.Error("encoder error:", err.Error())
	}

	return &buffer
}

func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
