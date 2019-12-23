package lib

import "github.com/astaxie/beego/utils/captcha"

var Cpt *captcha.Captcha

func init() {
	Cpt = captcha.NewCaptcha("", Cache)
	Cpt.ChallengeNums = 4
	Cpt.StdWidth = 100
	Cpt.StdHeight = 40
}
