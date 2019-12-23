package utils

type ResponseCode struct {
	Code int
	Msg  string
}

var RC map[string]ResponseCode

func init() {
	RC = map[string]ResponseCode{
		"success": ResponseCode{
			Code: 200,
			Msg:  "success",
		},
		"failed": ResponseCode{
			Code: 202,
			Msg:  "failed",
		},
		"invalid_args": {
			Code: 203,
			Msg:  "invalid arguments",
		},
		"redirect": {
			Code: 302,
			Msg:  "redirect",
		},
		"authentication_failed": {
			Code: 401,
			Msg: "authentication failed",
		},
	}
}

// 常用数据组装
func SimpleResAssembly(code int, msg string, data *map[string]interface{}) map[string]interface{} {
	if data == nil {
		data = &map[string]interface{}{}
	}
	return ResAssembly(code, msg, map[string]interface{}{"data": *data})
}

// response数据组装
func ResAssembly(code int, msg string, data ...map[string]interface{}) map[string]interface{} {
	r := map[string]interface{}{
		"code": code,
		"msg":  msg,
	}
	if data != nil {
		for _, v := range data {
			for k, item := range v {
				if item == nil {
					continue
				}
				r[k] = item
			}
		}
	}

	return r
}
