package utils

func SliceColumn(s []map[string]interface{}, column string) []interface{} {
	var res []interface{}
	for _, item := range s {
		if v, ok := item[column]; ok {
			res = append(res, v)
		}
	}
	return res
}