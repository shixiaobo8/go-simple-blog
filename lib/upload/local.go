package upload

import (
	"io"
	"os"
	"path"
	"strings"
)

type Local struct {
	baseDir string
}

func (l *Local) genearteDir() string {
	var p string
	var d = genDir()
	if l.baseDir != "" {
		p = path.Join(l.baseDir, d)
	}
	return p
}

func mkdir(dir string) error {
	exists, err := PathExists(dir)
	if err != nil {
		return err
	}
	if !exists {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Local) ReadFrom(src io.Reader, ext string) (string, error) {
	var err error
	var f *os.File

	dstDir := l.genearteDir()
	err = mkdir(dstDir)
	if err != nil {
		return "", err
	}
	fileName := genFileName(ext)
	dst := strings.TrimLeft(path.Join(dstDir, fileName), "\\/")
	f, err = os.OpenFile(dst, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, src)
	if err != nil {
		return "", err
	}
	return dst, nil
}

func (l *Local) FileUpload(file string) (string, error) {
	var err error
	src, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dstDir := l.genearteDir()
	err = mkdir(dstDir)
	if err != nil {
		return "", err
	}

	fileName := genFileName(path.Ext(file))
	dstFile := path.Join(dstDir, fileName)
	dst, err := os.OpenFile(dstFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", nil
	}
	return dstFile, nil
}


// Local
func NewLocalUploder(baseDir string) *Local {
	if baseDir != "" {
		baseDir = strings.TrimRight(baseDir, "\\/")
	}
	return &Local{baseDir:baseDir}
}

