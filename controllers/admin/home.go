package admin

import (
	"go-blog/controllers"
	"go-blog/utils"
)

type HomeController struct {
	controllers.AdminBaseController
}

func (this *HomeController) Index() {
	this.Data["json"] = utils.SimpleResAssembly(utils.RC["success"].Code, "success", nil)
	this.ServeJSON()
}
