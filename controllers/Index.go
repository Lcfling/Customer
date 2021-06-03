package controllers

import (
	//"github.com/virteman/OPMS/initial"

	"github.com/Lcfling/Customer/utils"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
)

type IndexController struct {
	beego.Controller
	IsLogin bool
	//UserInfo string
	Uid   int64
	UserUsername string
	UserAvatar   string
	UserAgent 	string
}

//hocker
func (this *IndexController) Prepare() {
	userLogin := this.GetSession("userLogin")
	useragent := this.Ctx.Request.Header.Get("user-agent")
	this.UserAgent=utils.GetUserAgent(useragent)
	if userLogin == nil {
		this.IsLogin = false
	} else {
		this.IsLogin = true
		tmp := strings.Split((this.GetSession("userLogin")).(string), "||")

		userid, _ := strconv.Atoi(tmp[0])
		longid := int64(userid)
		this.Data["LoginUserid"] = longid
		this.Data["LoginUsername"] = tmp[1]
		this.Data["LoginAvatar"] = tmp[2]

		this.Uid = longid
		this.UserUsername = tmp[1]
		this.UserAvatar = tmp[2]

		//this.Data["PermissionModel"] = this.GetSession("userPermissionModel")
		//this.Data["PermissionModelc"] = this.GetSession("userPermissionModelc")

	}
	this.Data["IsLogin"] = this.IsLogin
	//this.Data["IsLogin"] = this.IsLogin
}
func (this *IndexController) Options() {
	this.Data["json"] = map[string]interface{}{"status": 200, "message": "ok", "moreinfo": ""}
	this.ServeJSON()
}
