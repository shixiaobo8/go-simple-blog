package routers

import (
	"github.com/astaxie/beego"
	"go-blog/controllers/admin"
)

func init() {
	ns := beego.NewNamespace("admin",
		beego.NSRouter("/access/login", &admin.AccessController{}, "post:Login"),
		beego.NSRouter("/access/logout", &admin.AccessController{}, "post:Logout"),
		beego.NSRouter("/access/register", &admin.AccessController{}, "get:Register"),
		beego.NSRouter("/access/captcha/?:id", &admin.AccessController{}, "get:Captcha;post:Captcha"),

		beego.NSRouter("/home/index", &admin.HomeController{}, "get:Index"),

		beego.NSRouter("/adminUser/userInfo", &admin.AdminUserController{}, "get:UserInfo"),
		beego.NSRouter("/adminUser/list", &admin.AdminUserController{}, "post:List"),
		beego.NSRouter("/adminUser/updateUser", &admin.AdminUserController{}, "post:UpdateUser"),

		beego.NSRouter("/posts/list", &admin.PostsController{}, "post:List"),
		beego.NSRouter("/posts/info", &admin.PostsController{}, "get:Info"),
		beego.NSRouter("/posts/update", &admin.PostsController{}, "post:Update"),
		beego.NSRouter("/posts/delete", &admin.PostsController{}, "post:Delete"),
		beego.NSRouter("/posts/filterTag", &admin.PostsController{}, "post:FilterTag"),

		beego.NSRouter("/upload/image", &admin.UploadController{}, "post:Image"),
	)
	beego.AddNamespace(ns)
}
