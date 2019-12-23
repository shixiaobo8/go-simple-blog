package admin

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"go-blog/controllers"
	"go-blog/db"
	"go-blog/lib"
	"go-blog/models"
	"go-blog/utils"
	"net/url"
	"unicode/utf8"
)

type AdminUserController struct {
	controllers.AdminBaseController
}

func (this *AdminUserController) UserInfo() {
	var admin models.Admin
	var id, _ = this.GetInt("id", 0)

	if id == 0 {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["invalid_args"].Code, utils.RC["invalid_args"].Msg, nil)
		this.ServeJSON()
	}

	admin = models.Admin{}
	db.Db.Model(&models.Admin{}).Where("id = ?", id).First(&admin)

	if admin.ID == 0 {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["failed"].Code, "the user dose not exists", nil)
		this.ServeJSON()
	}

	avatar, _ := lib.Su.FormatPath(admin.Avatar)
	this.Data["json"] = utils.SimpleResAssembly(utils.RC["success"].Code, utils.RC["success"].Msg, &map[string]interface{}{
		"id":            admin.ID,
		"nickname":      admin.Nickname,
		"username":      admin.Username,
		"avatar":        avatar,
		"lastLoginTime": admin.LastLoginTime,
	})
	this.ServeJSON()
}

func (this *AdminUserController) List() {
	var m models.Admin
	var admins []models.Admin
	var res []map[string]interface{}

	columns := []map[string]interface{}{
		{
			"key" : "id",
			"title": "ID",
		},
		{
			"key" : "nickname",
			"title": "用户昵称",
		},
		{
			"key" : "avatar",
			"title": "头像",
		},
		{
			"key" : "last_login_time",
			"title" : "上次登录时间",
		},
	}

	fields := utils.SliceColumn(columns, "key")

	var fs []string
	for _, v := range fields {
		fs = append(fs, v.(string))
	}
	admins = m.List(fs)

	var avatar string
	for _, admin := range admins {
		avatar, _ = lib.Su.FormatPath(admin.Avatar)
		res = append(res, map[string]interface{}{
			"id" : admin.ID,
			"nickname" : admin.Nickname,
			"avatar" : avatar,
			"last_login_time" : admin.LastLoginTime,
		})
	}
	this.Data["json"] = utils.SimpleResAssembly(utils.RC["success"].Code, utils.RC["success"].Msg, &map[string]interface{}{
		"columns": columns,
		"rows" : res,
	})
	this.ServeJSON()
}

func (this *AdminUserController) UpdateUser() {
	var err error
	var ob = struct {
		Id uint
		Nickname string
		Avatar string
	}{}

	err = json.Unmarshal(this.Ctx.Input.RequestBody, &ob)
	if err != nil {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["invalid_args"].Code, utils.RC["invalid_args"].Msg, nil)
		this.ServeJSON()
	}

	if ob.Id <= 0 {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["invalid_args"].Code, utils.RC["invalid_args"].Msg, nil)
		this.ServeJSON()
	}

	if utf8.RuneCountInString(ob.Nickname) < 3 {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["invalid_args"].Code, "用户昵称不能少于3个字符", nil)
		this.ServeJSON()
	}

	// 将图片路径处理，仅存储部分路径
	if ob.Avatar != "" {
		urlParsed, err := url.Parse(ob.Avatar)
		if err != nil {
			this.Data["json"] = utils.SimpleResAssembly(utils.RC["invalid_args"].Code, "头像解析异常", nil)
			this.ServeJSON()
		}
		ob.Avatar = urlParsed.Path
	}

	_ = (&models.Admin{
		Model: gorm.Model{ID: ob.Id},
	}).UpdateAdminInfo(ob.Nickname, ob.Avatar)

	this.Data["json"] = utils.SimpleResAssembly(utils.RC["success"].Code, utils.RC["success"].Msg, nil)
	this.ServeJSON()
}
