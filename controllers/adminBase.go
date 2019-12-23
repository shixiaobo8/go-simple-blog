package controllers

import (
	"crypto/sha256"
	"github.com/astaxie/beego"
	"go-blog/models"
	"go-blog/utils"
	jwt2 "go-blog/utils/jwt"
	string2 "go-blog/utils/string"
	"regexp"
	"strings"
	"time"
)

type AdminBaseController struct {
	BaseController
}

var AdminSessionName = "admin"
type LoginSession struct {
	Id uint
	Nickname string
	UserName string
	Avatar string
	LastLoginTime *time.Time
}

func (this *AdminBaseController) jwtHeader() jwt2.Header {
	return jwt2.Header{
		AlgInstance: sha256.New(),
		Alg:         "sha256",
		Typ:         "jwt",
	}
}

func (this *AdminBaseController) Jwt() *jwt2.Jwt {
	header := this.jwtHeader()
	payload := jwt2.Payload{
		Iss: beego.AppConfig.String("appname"),
		Iat: time.Now(),
		Exp: time.Now().Add(2*time.Hour),
		Sub: this.Ctx.Request.Host,
		Jti: string2.GetRandomString(16),
	}

	extra := map[string]interface{}{}
	return jwt2.NewJwt(header, payload, extra)
}

// 登录设置session
func (this *AdminBaseController) SetLoginSession(admin models.Admin) {
	this.SetSession(AdminSessionName, LoginSession{
		Id:            admin.ID,
		Nickname:      admin.Nickname,
		UserName:      admin.Username,
		Avatar:        admin.Avatar,
		LastLoginTime: admin.LastLoginTime,
	})
}

// 验证用户登录状态可用session或jwt
func (this *AdminBaseController) LoginVerify() {
	// 用session验证
	//r := this.GetSession(AdminSessionName)
	//if r == nil {
	//	this.Ctx.ResponseWriter.WriteHeader(302)
	//	this.Data["json"] = utils.SimpleResAssembly(utils.RC["redirect"].Code, utils.RC["redirect"].Msg, &map[string]interface{}{
	//		"redirect_to": "/admin/access/login",
	//	})
	//	this.ServeJSON()
	//}

	// 用jwt验证
	jwt := jwt2.NewJwt(this.jwtHeader(), jwt2.Payload{}, map[string]interface{}{})
	token := this.Ctx.Request.Header.Get("Authorization")
	verify, err := jwt.VerifyToken(token)
	if err != nil {
		if e, ok := err.(jwt2.ExpireError); ok {
			if time.Now().Sub(e.Exp) > jwt2.RefreshSecond {
				goto ret
			} else {
				token, err := this.Jwt().Token()
				if err != nil {
					goto ret
				}
				this.Ctx.ResponseWriter.WriteHeader(401)
				this.Data["json"] = utils.SimpleResAssembly(utils.RC["authentication_failed"].Code, utils.RC["authentication_failed"].Msg, &map[string]interface{}{
					"token": token,
				})
				this.ServeJSON()
				return
			}
		}
		goto ret

	} else if verify == true {
		return
	}

	ret:
		this.Ctx.ResponseWriter.WriteHeader(401)
		this.Data["json"] = utils.SimpleResAssembly(utils.RC["authentication_failed"].Code, utils.RC["authentication_failed"].Msg, nil)
		this.ServeJSON()
}

func (this *AdminBaseController) Prepare() {
	match, _ := regexp.MatchString("^/admin/access/captcha.*", strings.ToLower(this.Ctx.Request.RequestURI))
	if match {
		return
	}

	if strings.ToLower(this.Ctx.Request.RequestURI) == "/admin/access/login" {
		return
	}
	this.LoginVerify()
}
