package upload

import (
	"github.com/rs/xid"
	"io"
	"os"
	"strconv"
	"time"
)

const (
	_ = iota
	LOCAL
	ALIYUN_OSS
)

type Uploader interface {
	genearteDir() string
	FileUpload(string) (string, error)
	ReadFrom(io.Reader, string) (string, error)
}

func genDir() string {
	y := strconv.Itoa(time.Now().Year())
	m := strconv.Itoa(int(time.Now().Month()))
	d := strconv.Itoa(time.Now().Day())
	return "/" + y + "/" + m + "/" + d
}

func genFileName(ext string) string {
	guid := xid.New()
	return guid.String() + ext
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
