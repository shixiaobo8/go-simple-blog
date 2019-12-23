package utils

import (
	"errors"
	"reflect"
)

// 获取结构体字段信息
func GetFieldNames(s interface{}) ([]string, error) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, errors.New("paramter one is not an struct")
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result, nil
}

// 判断结构体是否存在属性
func FieldExists(s interface{}, fieldName string) (bool, error) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false, errors.New("paramter one is not an struct")
	}
	if _, ok := t.FieldByName(fieldName); ok {
		return true, nil
	}
	return false, nil
}

// 获取结构体字段值
func GetFieldValue(s interface{}, fieldName string) (interface{}, error) {
	exists, err := FieldExists(s, fieldName)
	if err != nil {
		return nil, err
	}

	if exists == false {
		return nil, errors.New("struct has no field " + fieldName)
	}

	v := reflect.ValueOf(s)
	return v.FieldByName(fieldName).Interface(), nil
}

func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}