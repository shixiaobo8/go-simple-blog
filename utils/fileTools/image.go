package fileTools

import (
	"errors"
)

type fileTypeMap struct {
	fileType string
	header []byte
	typeLen int
}

var fileTypeSlice []fileTypeMap
var MaxTypeLen int

func init()  {
	fileTypeSlice = make([]fileTypeMap, 0)
	fileTypeSlice = append(fileTypeSlice, fileTypeMap{
		fileType: "jpg",
		header:   []byte{0xff, 0xd8},
		typeLen: 2,
	})

	fileTypeSlice = append(fileTypeSlice, fileTypeMap{
		fileType: "png",
		header:   []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
		typeLen: 8,
	})

	fileTypeSlice = append(fileTypeSlice, fileTypeMap{
		fileType: "gif",
		header:   []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61},
		typeLen: 6,
	})

	fileTypeSlice = append(fileTypeSlice, fileTypeMap{
		fileType: "gif",
		header:   []byte{0x47, 0x49, 0x46, 0x38, 0x37, 0x61},
		typeLen: 6,
	})

	for _, v := range fileTypeSlice {
		if v.typeLen > MaxTypeLen {
			MaxTypeLen = v.typeLen
		}
	}
}

// 检测文件类型
func ImageType(stream []byte) (string, error){
	var t string
	for _, v := range fileTypeSlice {
		streamLength := len(stream)
		length := len(v.header)

		if streamLength < length {
			continue
		}
		flag := true
		for i:=0; i<length; i++ {
			if stream[i] == v.header[i] {
				continue
			} else {
				flag = false
				break
			}
		}
		if flag {
			t = v.fileType
			break
		}
	}

	if t == "" {
		return t, errors.New("Undetectable type")
	}
	return t, nil
}