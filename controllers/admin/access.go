package admin

import (
	"encoding/json"
	"fmt"
	"go-blog/controllers"
	"go-blog/db"
	"go-blog/lib"
	"go-blog/models"
	"go-blog/utils"
	"time"
)

type AccessController struct {
	controllers.AdminBaseController
}

func (this *AccessController) Register() {
	// 添加默认账号
	var count = 0
	var err error

	if err = db.Db.Model(&models.Admin{}).Count(&count).Error; err == nil && count == 0 {
		db.Db.Create(&models.Admin{
			Username: "admin",
			Password: models.PasswordEnc("123456"),
			Nickname: "admin",
		})
	}
	this.Data["json"] = utils.SimpleResAssembly(utils.RC["success"].Code, "register success", nil)
	this.ServeJSON()
}

func (this *AccessController) Captcha() {
	if this.Ctx.Request.Method == "POST" {
		var err error
		var id string
		lib.Cpt.Expiration = 2 * 60 * time.Second
		id, err = lib.Cpt.CreateCaptcha()
		if err != nil {
			panic(err)
		}
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["success"].Code, utils.RC["success"].Msg, &map[string]interface{}{
			"id":  id,
			"cpt": fmt.Sprintf("%s", id),
		})
		this.ServeJSON()
	} else if this.Ctx.Request.Method == "GET" {
		lib.Cpt.Handler(this.Ctx)
	}
}

func (this *AccessController) Login() {
	var err error
	var admin models.Admin
	var ptrTime time.Time
	var token string
	var ob = struct {
		UserName string
		Password string
		CptId    string
		CptCode  string
	}{}
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &ob)

	if err != nil {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["invalid_args"].Code, utils.RC["invalid_args"].Msg, nil)
		goto End
	}
	if ob.CptCode == "" || ob.UserName == "" || ob.Password == "" {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["invalid_args"].Code, utils.RC["invalid_args"].Msg, nil)
		goto End
	}

	if lib.Cpt.Verify(ob.CptId, ob.CptCode) == false {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["invalid_args"].Code, "verification code error or invalid", nil)
		goto End
	}

	admin = models.Admin{}
	db.Db.Model(&models.Admin{}).Where("username = ? AND password = ?", ob.UserName, ob.Password).First(&admin)

	if admin.ID == 0 {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["failed"].Code, "the user dose not exists", nil)
		goto End
	}
	ptrTime = time.Now()
	admin.UpdateLastLoginTime(&ptrTime)
	this.SetLoginSession(admin)
	token, err = this.Jwt().Token()
	if err != nil {
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["failed"].Code, "token issuance of failure", nil)
		goto End
	}
	this.Data["json"] = utils.SimpleResAssembly(utils.RC["success"].Code, "login success", &map[string]interface{}{
		"id": admin.ID,
		"token": token,
	})
End:
	this.ServeJSON()
}

func (this *AccessController) Logout () {
	this.DestroySession()
	this.Data["json"] = utils.SimpleResAssembly(utils.RC["success"].Code, "logout", nil)
	this.ServeJSON()
}
