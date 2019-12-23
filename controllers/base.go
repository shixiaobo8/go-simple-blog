package controllers

import (
	"github.com/astaxie/beego"
	"net/http"
)

type BaseController struct {
	beego.Controller
}

// 获取sessionid
func (this *BaseController) GetSid() (string, error) {
	var cookie *http.Cookie
	var err error
	cookie, err = this.Ctx.Request.Cookie(beego.BConfig.WebConfig.Session.SessionName)
	if err != nil {
		return "", err
	}
	return cookie.Value, err
}
