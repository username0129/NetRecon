package util

import (
	"reflect"
)

// StructToMap
//
//	@Description: 将结构体转换为 map 类型
//	@param obj
//	@return map[string]interface{}
func StructToMap(obj interface{}) map[string]interface{} {
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)

	// 处理指针到结构体的情况
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
		objValue = objValue.Elem()
	}

	// 确保是结构体
	if objType.Kind() != reflect.Struct {
		return nil
	}

	data := make(map[string]interface{})
	// 遍历字段
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldValue := objValue.Field(i)
		// 跳过非导出字段
		if !fieldValue.CanInterface() {
			continue
		}
		key := field.Name
		if tag := field.Tag.Get("mapstructure"); tag != "" {
			key = tag
		}
		data[key] = fieldValue.Interface()
	}
	return data
}
