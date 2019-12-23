package main

import (
	"encoding/gob"
	"go-blog/controllers"
	_ "go-blog/routers"
	_ "go-blog/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/astaxie/beego/session/redis"
)

func main() {
	// 允许跨域
	// Access-Control-Allow-Credentials设置为true的情况下 Access-Control-Allow-Origin不能设置为*
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:8080", "http://127.0.0.1:8080"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
	}))

	initSession()
	beego.Run()
}

func initSession() {
	gob.Register(controllers.LoginSession{})
}
