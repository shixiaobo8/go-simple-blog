package routers

import (
	"github.com/astaxie/beego"
	"go-blog/controllers/api"
)

func init() {
	ns := beego.NewNamespace("api",
		beego.NSRouter("/posts/list", &api.PostsController{}, "post:List"),
		beego.NSRouter("/posts/info", &api.PostsController{}, "get:Info"),
		beego.NSRouter("/posts/recommend", &api.PostsController{}, "post:Recommend"),
	)
	beego.AddNamespace(ns)
}
