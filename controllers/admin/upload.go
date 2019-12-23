package admin

import (
	"go-blog/controllers"
	"go-blog/lib"
	"go-blog/utils"
	"go-blog/utils/fileTools"
	"io"
)

type UploadController struct {
	controllers.AdminBaseController
}

// 单文件上传
func (this *UploadController) Image() {
	file, _, err := this.GetFile("file")
	defer file.Close()

	var fh []byte
	fh = make([]byte, fileTools.MaxTypeLen)
	_, err = file.Read(fh)
	if err != nil  && err != io.EOF {
		this.Data["json"] =  utils.SimpleResAssembly(utils.RC["failed"].Code, "图片上传失败", nil)
		this.ServeJSON()
	}
	_, err = file.Seek(0,0)
	if err != nil  && err != io.EOF {
		this.Data["json"] =  utils.SimpleResAssembly(utils.RC["failed"].Code, "图片上传失败", nil)
		this.ServeJSON()
	}

	ext, err := fileTools.ImageType(fh)
	if err != nil {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["failed"].Code, "不支持的图片格式", nil)
		this.ServeJSON()
	}
	ext = "." + ext

	if err != nil {
		this.Data["json"] =  utils.SimpleResAssembly(utils.RC["failed"].Code, "图片上传失败", nil)
		this.ServeJSON()
	}

	var img string
	img, err = lib.Su.Uploader.ReadFrom(file, ext)
	if err != nil {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["failed"].Code, "图片上传失败", nil)
		this.ServeJSON()
	}

	img, err = lib.Su.FormatPath(img)
	if err != nil {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["failed"].Code, "图片已上传但处理失败", nil)
		this.ServeJSON()
	}

	this.Data["json"] = utils.SimpleResAssembly(utils.RC["success"].Code, "success", &map[string]interface{}{
		"img" : img,
	})
	this.ServeJSON()
}
