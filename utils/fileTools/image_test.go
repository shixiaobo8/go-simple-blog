package fileTools

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestImageType(t *testing.T) {
	var stream []byte = []byte{137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 120}
	logs.Info(ImageType(stream))
}