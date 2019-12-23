package models

import "database/sql"

/**
常用功能函数
 */


// 将Row转换成map切片
func RowsToMap(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	var list []map[string]interface{}
	var r  = make([]interface{}, len(columns))	// 存储每行数据

	for i, _ := range r {
		var a interface{}
		r[i] = &a		// 利用指针，Scan使用
	}

	for rows.Next() {
		_ = rows.Scan(r...)
		item := make(map[string]interface{}) // 再转换成map
		for i, v := range r {
			item[columns[i]] = *v.(*interface{})
		}
		list = append(list, item)
	}
	return list
}
