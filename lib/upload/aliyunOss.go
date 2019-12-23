package upload

import (
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"os"
	"path"
	"strings"
)

type AliyunOss struct {
	baseDir string
	bucket *oss.Bucket
	client *oss.Client
}

func (a *AliyunOss) genearteDir() string {
	var p string
	var d = genDir()
	if a.baseDir != "" {
		p = path.Join(a.baseDir, d)
	} else {
		p = d
	}
	return p
}

func (a *AliyunOss) getBucket() *oss.Bucket {
	return a.bucket
}

func (a *AliyunOss) getClient() *oss.Client {
	return a.client
}

func (a *AliyunOss) SetBucket(bucketName string) error {
	bucket, err := a.client.Bucket(bucketName)
	if err != nil {
		return err
	}
	a.bucket = bucket
	return nil
}

func (a *AliyunOss) ReadFrom(src io.Reader, ext string) (string, error) {
	var err error
	dstDir := a.genearteDir()
	fileName := genFileName(ext)
	dst := strings.TrimLeft(path.Join(dstDir, fileName), "\\/")
	if a.bucket == nil {
		return "", errors.New("property Bucket is null")
	}
	err = a.bucket.PutObject(dst, src)
	if err != nil {
		return "", err
	}
	return dst, nil
}

func (a *AliyunOss) FileUpload(file string) (string, error) {
	var err error
	if _, err = os.Lstat(file); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("文件不存在")
		}
	}
	dstDir := a.genearteDir()
	fileName := genFileName(path.Ext(file))
	dst := strings.TrimLeft(path.Join(dstDir, fileName), "\\/")
	if a.bucket == nil {
		return "", errors.New("property Bucket is null")
	}
	err = a.bucket.PutObjectFromFile(dst, file)
	if err != nil {
		return "", err
	}
	return dst, nil
}

// AliyunOss
func NewAliyunOssUploader(baseDir string, endPoint string, accessKeyId string, accessKeySecret string) (*AliyunOss, error) {
	if baseDir != "" {
		baseDir = strings.TrimRight(baseDir, "\\/")
	}

	client, err := oss.New(endPoint, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, err
	}

	return &AliyunOss{
		baseDir: baseDir,
		bucket:  nil,
		client:  client,
	}, nil
}
