package util

import (
	"errors"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"reflect"
	"strconv"
	"strings"
)

// GenerateStructField 生成struct属性
func GenerateStructField(metaDatas []*ConfigMetaData) string {
	var build strings.Builder
	for _, metaData := range metaDatas {
		//首字母大写，否则不能反射
		build.WriteString(util.FirstUpper(metaData.Name))
		build.WriteString(" ")
		build.WriteString(metaData.Type)
		build.WriteString(" ")
		//`json:"id"`
		build.WriteString("`json:\"")
		build.WriteString(metaData.Name)
		build.WriteString("\"` //")
		build.WriteString(metaData.Description)
		build.WriteString("\n")
	}

	return build.String()
}

// UnmarshalConfig 解析配置文件,只支持一层,不支持slice，map
func UnmarshalConfig(metaDatas []*ConfigMetaData, v interface{}) error {
	if metaDatas == nil {
		return nil
	}
	vType := reflect.TypeOf(v)
	vValue := reflect.ValueOf(v)
	if vType.Kind() != reflect.Ptr { //因为要修改v，必须传指针
		return errors.New("must pass pointer parameter")
	}
	vType = vType.Elem()
	vValue = vValue.Elem()
	if vType.Kind() != reflect.Struct {
		return errors.New("value must struct")
	}

	fieldCount := vType.NumField()
	tagNames := make(map[string]string, fieldCount)
	for i := 0; i < fieldCount; i++ {
		fieldType := vType.Field(i)
		name := fieldType.Name
		if len(fieldType.Tag.Get("json")) > 0 {
			name = fieldType.Tag.Get("json")
		}
		tagNames[name] = fieldType.Name
	}
	for _, metaData := range metaDatas {
		if name, exist := tagNames[metaData.Name]; exist {
			value := vValue.FieldByName(name)
			//不支持指针
			if value.Kind() == reflect.Ptr {
				continue
			}
			switch metaData.Type {
			case "int",
				"int8",
				"int16",
				"int32",
				"int64":
				if i, err := strconv.ParseInt(metaData.Value, 10, 64); err != nil {
					return err
				} else {
					value.SetInt(i) //有符号整型通过SetInt
				}

			case "uint",
				"uint8",
				"uint16",
				"uint32",
				"uint64":
				if i, err := strconv.ParseUint(metaData.Value, 10, 64); err != nil {
					return err
				} else {
					value.SetUint(i) //无符号整型需要通过SetUint
				}

			case "string":
				value.SetString(metaData.Value)
			case "bool":
				if b, err := strconv.ParseBool(metaData.Value); err == nil {
					value.SetBool(b)
				} else {
					return err
				}
			case "float32",
				"float64":
				if f, err := strconv.ParseFloat(metaData.Value, 64); err != nil {
					return err
				} else {
					value.SetFloat(f) //通过reflect.Value修改原始数据的值
				}
			default:
				log.Warn("config [%v] not set value,config error?", metaData)
			}
		} else {
			log.Debug("struct field [%v] not define", metaData.Name)
			continue
		}

	}

	return nil
}

// ConfigMetaData 配置元数据
type ConfigMetaData struct {
	Name        string `json:"name"`        //属性名称
	Value       string `json:"value"`       //属性值
	Type        string `json:"type"`        //属性类型
	Description string `json:"description"` //属性描述
}
