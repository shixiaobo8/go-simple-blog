package lib

import (
	"github.com/astaxie/beego"
	"go-blog/lib/upload"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var Su SimpleUpload

type SimpleUpload struct {
	upType int
	Uploader upload.Uploader
}

func init() {
	var err error
	Su.upType, err = beego.AppConfig.Int("uploader")
	if err != nil {
		panic(err)
	}

	if Su.upType == upload.ALIYUN_OSS {
		endpoint := beego.AppConfig.String("aliyunoss_endpoint")
		accesskeyid := beego.AppConfig.String("aliyunoss_accesskeyid")
		accesskeysecret := beego.AppConfig.String("aliyunoss_accesskeysecret")
		bucket := beego.AppConfig.String("aliyunoss_default_bucket")

		Su.Uploader, err = upload.NewAliyunOssUploader("", endpoint, accesskeyid, accesskeysecret)
		u := Su.Uploader.(*upload.AliyunOss)
		if err != nil {
			panic(err)
		}
		err = u.SetBucket(bucket)
		if err != nil {
			panic(err)
		}
	} else {
		uploadpath := beego.AppConfig.String("uploadpath")
		if !filepath.IsAbs(uploadpath) {
			uploadpath, err = filepath.Abs(uploadpath)
			if err != nil {
				panic(err)
			}
		}
		Su.Uploader = upload.NewLocalUploder(uploadpath)
	}
}

// 本地上传后文件名为本地绝对路径，oss文件上传后文件名为oss的绝对路径
// 将路径进行转换成域名形式供前端使用
func (su SimpleUpload) FormatPath(file string) (string, error) {
	var domain string
	if su.upType == upload.ALIYUN_OSS {
		domain = beego.AppConfig.String("aliyunoss_default_bucket_domain")
		return strings.TrimRight(domain, "/") + "/" + file, nil
	} else {
		domain = beego.AppConfig.String("httpscheme") + "://" + beego.AppConfig.String("httpaddr")
		pwd, _ := os.Getwd()
		p, err := filepath.Rel(pwd, file)
		if err != nil {
			return "", err
		}
		port, err := beego.AppConfig.Int("httpport")
		if err != nil {
			return "", err
		}
		if port != 0 {
			return strings.Replace(domain + ":" + strconv.Itoa(port) + "/" + strings.TrimRight(p, "./"), "\\", "/", -1), nil
		} else {
			return strings.Replace(domain + "/" + strings.TrimRight(p, "./"), "\\", "/", -1), nil
		}
	}
}
