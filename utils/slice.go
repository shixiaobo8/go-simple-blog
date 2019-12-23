package utils

import (
	"errors"
	"reflect"
)

// 去重
func SliceRemoveDupString(s []string) []string {
	result := []string{}
	tempMap := map[string]byte{}
	for _, e := range s {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

// 获取切片结构体中的某个字段值
func SliceStructColumn(v []interface{}, property string) ([]interface{}, error) {
	var res []interface{}
	for _, item := range v {
		t := reflect.TypeOf(item)
		if t.Kind() != reflect.Struct {
			return res, errors.New("invalid paramter")
		}
		r, err := GetFieldValue(item, property)
		if err != nil {
			return res, err
		}
		res = append(res, r)
	}
	return res, nil
}

// 切片交集
func SliceIntersection(s1 []interface{}, s2 []interface{}) []interface{} {
	var min []interface{}
	var max []interface{}

	if len(s1) >= len(s2) {
		min = s2
		max = s1
	} else {
		min = s1
		max = s2
	}

	var res []interface{}
	for _, item := range min {
		for k, value := range max {
			if item == value {
				res = append(res, item)
				max = append(max[0:k], max[k+1:]...)
				break
			}
		}
	}
	return res
}

// 切片差集
func SliceDifference(s1 []interface{}, s2 []interface{}) []interface{} {
	var res []interface{}
	for _, item := range s1 {
		exists := false
		for _, value := range s2 {
			if item == value {
				exists = true
				break
			}
		}
		if exists == false {
			res = append(res, item)
		}
	}
	return res
}